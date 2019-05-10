const app = require('./app')
const config = require('./config')
const port = config.PORT

app.listen(port, () => console.log(`Example app listening on port ${port}!`))
