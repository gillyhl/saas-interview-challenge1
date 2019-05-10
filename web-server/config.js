const path = require('path')
module.exports = {
  PORT: process.env.PORT || 5000,
  UPLOAD_LOCATION:
    process.env.UPLOAD_LOCATION || path.join(__dirname, 'uploads'),
  REDIS_HOST: process.env.REDIS_HOST || 'localhost'
}
