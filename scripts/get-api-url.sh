#!/bin/bash

echo "Fetching API URL..."

TASK_ARN=$(aws ecs list-tasks --cluster book-boy-cluster --service-name book-boy-service --query "taskArns[0]" --output text --region us-east-2)

if [ "$TASK_ARN" == "None" ] || [ -z "$TASK_ARN" ]; then
  echo "Error: No running tasks found"
  exit 1
fi

ENI=$(aws ecs describe-tasks --cluster book-boy-cluster --tasks $TASK_ARN --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" --output text --region us-east-2)

if [ -z "$ENI" ]; then
  echo "Error: Could not find network interface"
  exit 1
fi

IP=$(aws ec2 describe-network-interfaces --network-interface-ids $ENI --query "NetworkInterfaces[0].Association.PublicIp" --output text --region us-east-2)

if [ -z "$IP" ]; then
  echo "Error: Could not find public IP"
  exit 1
fi

echo ""
echo "======================================"
echo "API URL: http://$IP:8080"
echo "======================================"
echo ""
echo "Test with:"
echo "curl http://$IP:8080/auth/register -X POST -H \"Content-Type: application/json\" -d '{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}'"
