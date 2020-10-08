const Host = process.env.DBHOST
const User = process.env.SQLUSER
const Pass = process.env.SQLPASS

const knex = require('knex')({
  client: 'mysql',
  connection: {
    host : Host,
    user : User,
    password : Pass,
    database : "Vtuber"
  }
});

module.exports = knex;