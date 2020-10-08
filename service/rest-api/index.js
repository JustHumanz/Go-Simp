const express = require("express");
const bodyParser = require("body-parser");

const app = express();

// parse requests of content-type: application/json
app.use(bodyParser.json());

// parse requests of content-type: application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: true }));

// simple route
app.get("/", (req, res) => {
  res.json({ message: "https://github.com/JustHumanz/Go-simp/blob/master/service/rest-api/routes/routes.js" });
});

require("./routes/routes.js")(app);

// set port, listen for requests
const PORT = process.env.PORT || 2525;

app.listen(PORT, () => {
  console.log("Server is running on port "+PORT);
});
