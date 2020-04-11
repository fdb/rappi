const express = require("express");
const path = require("path");
const Twit = require("twit");

require("dotenv").config();

const T = new Twit({
  consumer_key: process.env.TWITTER_CONSUMER_KEY,
  consumer_secret: process.env.TWITTER_CONSUMER_SECRET,
  access_token: process.env.TWITTER_ACCESS_TOKEN,
  access_token_secret: process.env.TWITTER_ACCESS_TOKEN_SECRET
});

const STATIC_DIR = path.join(__dirname, "static");

const app = express();
app.use(express.static(STATIC_DIR));
app.engine('html', require('ejs').renderFile);

app.get("/", (req, res) => res.render("index.html"));
app.get("/twitter", (req, res) => res.render("twitter.html"));
app.get("/twitter/search.json", async (req, res) => {
  console.log(req.params);
  try {
    const data = await T.get('search/tweets', { q: req.query.q, count: 100 });
    res.send({ status: 'ok', data });
  } catch (err) {
    res.send({ status: 'err', message: err });
  }
});

const port = process.env.PORT || 3000;
app.listen(port, () =>
  console.log(`Rappi listening at http://localhost:${port}`)
);
