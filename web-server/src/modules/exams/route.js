const config = require('../../../config')
const multer = require('multer')
const upload = multer({
  dest: config.UPLOAD_LOCATION
})
const fs = require('fs')
const uniqid = require('uniqid')
const service = require('./service')
const apiUtils = require('../../services/apiUtils')
module.exports = (app, redisClient) => {
  const router = require('express').Router()
  app.use('/exams', router)

  /**
   * Post an exam file to redis to be marked.
   */
  router.post('/', upload.single('exam'), (req, res) => {
    const id = uniqid()
    const contents = fs.readFileSync(req.file.path, { encoding: 'utf-8' })
    const jsonContents = {
      id,
      ...JSON.parse(contents)
    }
    redisClient.publish('exam', JSON.stringify(jsonContents))
    return apiUtils.handleSuccess(res, { id }, 201)
  })

  /**
   * Get the marked papers for the exam given id.
   */
  router.get('/:id', async (req, res) => {
    const { id } = req.params
    const { data, error } = await apiUtils.to(
      service.getMarkedPapersByExamId(id, redisClient)
    )
    if (error)
      return apiUtils.handleError(res, 'Unable to get exam marked papers', 500)
    return apiUtils.handleSuccess(res, data)
  })

  /**
   * Get the marked paper given exam id and paper id
   */
  router.get('/:examId/papers/:paperId', async (req, res) => {
    const { examId, paperId } = req.params
    const { data, error } = await apiUtils.to(
      service.getMarkedPaperByExamIdAndPaperId(examId, paperId, redisClient)
    )
    if (error)
      return apiUtils.handleError(res, 'Unable to get exam marked paper', 500)
    return apiUtils.handleSuccess(res, data)
  })
}
