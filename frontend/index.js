const fs = require('fs');
const mysql = require('mysql');

// Create a MySQL connection
const connection = mysql.createConnection({
  host: '127.0.0.1',
  user: 'newuser',
  password: 'password',
  database: 'multicloud'
});

connection.connect((err) => {
    if (err) throw err;
    console.log('Connected to the database');
  });
  
  // Read the file content synchronously
  const filePath = '/home/yuyu/Desktop/divine-treat-413716-8ae024729d4f.json'; // Replace with the path to your file
  const fileContent = fs.readFileSync(filePath, { encoding: 'base64' });
  
  // Prepare a SQL query to update the file content in the database
  const sql = 'UPDATE CloudAccount SET KeyFile = ? WHERE AccountID = 5';
  
  // Execute the SQL query with the file content and AccountID as parameters
  connection.query(sql, [fileContent, 5], (err, result) => {
    if (err) throw err;
    console.log('File content updated in the database');
  });
  
  // Close the database connection
  connection.end();