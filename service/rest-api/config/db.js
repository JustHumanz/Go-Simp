const Host = process.env.DBHOST
const User = process.env.DBUSER
const Pass = process.env.DBPASS

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