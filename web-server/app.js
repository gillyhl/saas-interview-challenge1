const express = require('express')
const config = require('./config')
const glob = require('glob')
const path = require('path')
const app = express()
const redis = require('redis')

const router = express.Router()
app.use('/api/v1', router)
const redisClient = redis.createClient({
  host: config.REDIS_HOST
})
const matches = glob.sync(path.join(__dirname, './src/modules/**/route.js'))
matches.forEach(match => {
  require(match)(router, redisClient)
})

app.get('/', (req, res) => res.send('Hello World!'))

module.exports = app
