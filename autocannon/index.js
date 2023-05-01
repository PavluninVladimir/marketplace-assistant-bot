const autocannon = require('autocannon');

const opts = {
  url: 'http://localhost:8181',
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