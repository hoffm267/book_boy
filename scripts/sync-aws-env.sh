#!/bin/bash

echo "WARNING: Run after deployment completes (5-10 minutes after git push)"

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
echo "API URL: http://$IP:8080"
echo ""

BRUNO_FILE="bruno/environments/aws.bru"

if [ ! -f "$BRUNO_FILE" ]; then
  echo "Error: Bruno AWS environment file not found at $BRUNO_FILE"
  exit 1
fi

sed -i "s|baseUrl: http://[0-9.]*:8080|baseUrl: http://$IP:8080|g" "$BRUNO_FILE"

echo "Bruno AWS environment updated!"
echo ""
