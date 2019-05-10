const { promisify } = require('util')
const getMarkedPapersByExamId = async (examId, redisClient) => {
  const getAsync = promisify(redisClient.lrange).bind(redisClient)
  const key = `exam#${examId}:marked-papers`
  const data = await getAsync(key, 0, -1)
  return data.map(d => JSON.parse(d))
}

const getMarkedPaperByExamIdAndPaperId = async (
  examId,
  paperId,
  redisClient
) => {
  const getAsync = promisify(redisClient.get).bind(redisClient)
  const key = `exam#${examId}:marked-papers#${paperId}`
  const data = await getAsync(key)
  return JSON.parse(data)
}

module.exports = {
  getMarkedPapersByExamId,
  getMarkedPaperByExamIdAndPaperId
}
