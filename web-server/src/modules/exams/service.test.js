const service = require('./service')
describe('exams service', () => {
  it('should get marked paper by exam id', async () => {
    const jsonString = ['{"a":1}', '{"a":2}']
    const examId = 'exam1'
    const lrangeStub = sinon.stub().callsArgWith(3, null, jsonString)
    const redisClient = {
      lrange: lrangeStub
    }

    const actualResult = await service.getMarkedPapersByExamId(
      examId,
      redisClient
    )
    const expectedResult = [{ a: 1 }, { a: 2 }]
    const stubCallArguments = lrangeStub.getCalls()[0].args
    expect(expectedResult).to.deep.equal(actualResult)
    expect('exam#exam1:marked-papers').to.deep.equal(stubCallArguments[0])
  })

  it('should get marked paper by exam and paper id', async () => {
    const jsonString = '{"a":1}'
    const examId = 'exam1'
    const paperId = 'paper1'
    const getStub = sinon.stub().callsArgWith(1, null, jsonString)
    const redisClient = {
      get: getStub
    }

    const actualResult = await service.getMarkedPaperByExamIdAndPaperId(
      examId,
      paperId,
      redisClient
    )
    const expectedResult = { a: 1 }
    const stubCallArguments = getStub.getCalls()[0].args
    expect(expectedResult).to.deep.equal(actualResult)
    expect('exam#exam1:marked-papers#paper1').to.deep.equal(
      stubCallArguments[0]
    )
  })
})
