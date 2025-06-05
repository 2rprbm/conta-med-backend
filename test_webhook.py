import requests
import json

# Test webhook verification endpoint
def test_webhook_verification():
    url = "http://44.204.27.5:8080/webhook"
    params = {
        "hub.mode": "subscribe",
        "hub.verify_token": "contamed_webhook_2024_secure",
        "hub.challenge": "test123"
    }
    
    try:
        response = requests.get(url, params=params, timeout=10)
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.text}")
        return response.status_code == 200
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return False

# Test webhook POST endpoint
def test_webhook_post():
    url = "http://44.204.27.5:8080/webhook"
    payload = {
        "object": "whatsapp_business_account",
        "entry": [
            {
                "id": "0",
                "changes": [
                    {
                        "field": "messages",
                        "value": {
                            "messaging_product": "whatsapp",
                            "metadata": {
                                "display_phone_number": "16505551111",
                                "phone_number_id": "123456123"
                            },
                            "contacts": [
                                {
                                    "profile": {
                                        "name": "test user name"
                                    },
                                    "wa_id": "16315551181"
                                }
                            ],
                            "messages": [
                                {
                                    "from": "16315551181",
                                    "id": "ABGGFlA5Fpa",
                                    "timestamp": "1504902988",
                                    "type": "text",
                                    "text": {
                                        "body": "this is a test message"
                                    }
                                }
                            ]
                        }
                    }
                ]
            }
        ]
    }
    
    try:
        response = requests.post(url, json=payload, timeout=10)
        print(f"POST Status Code: {response.status_code}")
        print(f"POST Response: {response.text}")
        return response.status_code == 200
    except requests.exceptions.RequestException as e:
        print(f"POST Error: {e}")
        return False

# Test basic connectivity
def test_connectivity():
    url = "http://44.204.27.5:8080/ping"
    try:
        response = requests.get(url, timeout=5)
        print(f"Ping Status: {response.status_code}")
        print(f"Ping Response: {response.text}")
        return True
    except requests.exceptions.RequestException as e:
        print(f"Connectivity Error: {e}")
        return False

if __name__ == "__main__":
    print("=== Testing AWS Application ===")
    print(f"Testing IP: 44.204.27.5:8080")
    print()
    
    print("1. Testing basic connectivity...")
    test_connectivity()
    print()
    
    print("2. Testing webhook verification (GET /webhook)...")
    test_webhook_verification()
    print()
    
    print("3. Testing webhook POST...")
    test_webhook_post()
    print()
    
    print("=== Test Complete ===")
    print()
    print("🔧 NEXT STEPS:")
    print("   1. Deploy the updated code to AWS")
    print("   2. Set up HTTPS with ALB + SSL Certificate")
    print("   3. Update Meta Developer with HTTPS URL")
    print()
    print("📍 CURRENT URLs:")
    print("   HTTP: http://44.204.27.5:8080/webhook")
    print("   HTTPS: [Need to set up ALB first]")
    print("   Verify Token: contamed_webhook_2024_secure") 