module.exports = {
  to: promise =>
    promise
      .then(data => ({
        data
      }))
      .catch(error => ({
        error
      })),
  handleSuccess: (res, data, statusCode = 200) => {
    res.status(statusCode).json({
      data,
      success: true,
      statusCode
    })
  },
  handleError: (res, message, statusCode = 500) => {
    res.status(statusCode).json({
      data: {
        reason: message
      },
      success: false,
      statusCode
    })
  }
}
