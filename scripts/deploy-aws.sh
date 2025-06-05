#!/bin/bash

# ContaMed Backend - AWS Deploy Script
# This script deploys the application to AWS using ECR and ECS

set -e

# Configuration
REGION="${AWS_REGION:-us-east-1}"
PROJECT_NAME="conta-med-backend"
ECR_REPO_NAME="conta-med-backend"
ECS_CLUSTER_NAME="conta-med-cluster"
ECS_SERVICE_NAME="conta-med-service"
TASK_DEFINITION_NAME="conta-med-task"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if AWS CLI is installed
check_aws_cli() {
    if ! command -v aws &> /dev/null; then
        log_error "AWS CLI is not installed. Please install it first."
        exit 1
    fi
    
    # Check if AWS is configured
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS CLI is not configured. Run 'aws configure' first."
        exit 1
    fi
    
    log_info "AWS CLI is configured correctly"
}

# Get AWS Account ID
get_account_id() {
    ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
    log_info "Using AWS Account ID: $ACCOUNT_ID"
}

# Create ECR repository if it doesn't exist
create_ecr_repo() {
    log_info "Creating ECR repository if it doesn't exist..."
    
    if ! aws ecr describe-repositories --repository-names $ECR_REPO_NAME --region $REGION &> /dev/null; then
        aws ecr create-repository \
            --repository-name $ECR_REPO_NAME \
            --region $REGION \
            --image-scanning-configuration scanOnPush=true
        log_info "ECR repository created: $ECR_REPO_NAME"
    else
        log_info "ECR repository already exists: $ECR_REPO_NAME"
    fi
}

# Build and push Docker image
build_and_push_image() {
    log_info "Building and pushing Docker image..."
    
    # Get ECR login token
    aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com
    
    # Build image
    IMAGE_URI="$ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/$ECR_REPO_NAME:latest"
    log_info "Building image: $IMAGE_URI"
    
    docker build -t $ECR_REPO_NAME .
    docker tag $ECR_REPO_NAME:latest $IMAGE_URI
    
    # Push image
    docker push $IMAGE_URI
    log_info "Image pushed successfully: $IMAGE_URI"
}

# Create ECS cluster if it doesn't exist
create_ecs_cluster() {
    log_info "Creating ECS cluster if it doesn't exist..."
    
    if ! aws ecs describe-clusters --clusters $ECS_CLUSTER_NAME --region $REGION --query 'clusters[0].status' --output text 2>/dev/null | grep -q ACTIVE; then
        aws ecs create-cluster \
            --cluster-name $ECS_CLUSTER_NAME \
            --region $REGION \
            --capacity-providers FARGATE \
            --default-capacity-provider-strategy capacityProvider=FARGATE,weight=1
        log_info "ECS cluster created: $ECS_CLUSTER_NAME"
    else
        log_info "ECS cluster already exists: $ECS_CLUSTER_NAME"
    fi
}

# Main deployment function
deploy() {
    log_info "Starting deployment to AWS..."
    log_info "Region: $REGION"
    log_info "Project: $PROJECT_NAME"
    
    check_aws_cli
    get_account_id
    create_ecr_repo
    build_and_push_image
    create_ecs_cluster
    
    log_info "Deployment completed successfully!"
    log_info "Image URI: $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/$ECR_REPO_NAME:latest"
    log_warn "Next steps:"
    log_warn "1. Create ECS Task Definition"
    log_warn "2. Create ECS Service"
    log_warn "3. Configure Load Balancer"
    log_warn "4. Set up domain and SSL"
}

# Run deployment
deploy 