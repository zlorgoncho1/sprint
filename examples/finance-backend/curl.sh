curl -X POST http://localhost:8000/finance/addIncome -H "Content-Type: application/json" -d '{"amount": 5000, "category": "Salary"}'

curl -X POST http://localhost:8000/finance/addExpense -H "Content-Type: application/json" -d '{"amount": 2000, "category": "Food"}'

curl -X GET http://localhost:8000/finance/getSummary

curl -X POST http://localhost:8000/user/register -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpass"}'

curl -X POST http://localhost:8000/user/login -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpass"}'

curl -X GET http://localhost:8000/user/profile\?username\=testuser
