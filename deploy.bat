@echo off
echo ========================================
echo Deploying conta-med-backend
echo ========================================

set ECR_REPOSITORY=323195246434.dkr.ecr.us-east-1.amazonaws.com/conta-med-backend
set TAG=latest
set REGION=us-east-1
set CLUSTER=conta-med-cluster
set SERVICE=conta-med-service

echo.
echo 1. Fazendo login no ECR...
aws ecr get-login-password --region %REGION% | docker login --username AWS --password-stdin %ECR_REPOSITORY%

echo.
echo 2. Construindo imagem Docker...
docker build -t conta-med-backend:%TAG% .

echo.
echo 3. Tageando imagem...
docker tag conta-med-backend:%TAG% %ECR_REPOSITORY%:%TAG%

echo.
echo 4. Enviando imagem para o ECR...
docker push %ECR_REPOSITORY%:%TAG%

echo.
echo 5. Atualizando serviço ECS...
aws ecs update-service --cluster %CLUSTER% --service %SERVICE% --force-new-deployment --region %REGION%

echo.
echo ========================================
echo IMPLANTAÇÃO INICIADA COM SUCESSO!
echo ========================================
echo A atualização do serviço foi iniciada.
echo A implantação pode levar alguns minutos.
echo.
echo Para verificar o status:
echo 1. Acesse: https://console.aws.amazon.com/ecs/
echo 2. Selecione o cluster: %CLUSTER%
echo 3. Vá para o serviço: %SERVICE%
echo.
echo Após a implantação, teste acessando:
echo http://conta-med-alb-501398913.us-east-1.elb.amazonaws.com/webhook
echo ========================================

pause 