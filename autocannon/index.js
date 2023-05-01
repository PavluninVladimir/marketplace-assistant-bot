const autocannon = require('autocannon');

const opts = {
  url: 'https://bot.my-infant.com/static/',
  connections: 100,
  duration: 10
};

autocannon(opts, (err, result) => {
  if (err) {
    console.error(err);
  } else {
    console.log(result);
  }
});