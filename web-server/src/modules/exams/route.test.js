const redis = require('redis')
const request = require('supertest')
const service = require('./service')

describe('exams route', function() {
  before(() => {
    this.publishStub = sinon.stub().returns()
    this.client = {
      publish: this.publishStub
    }
    this.createClientStub = sinon
      .stub(redis, 'createClient')
      .returns(this.client)
  })

  after(() => {
    this.createClientStub.restore()
  })

  beforeEach(() => {
    this.sandbox = sinon.createSandbox()
  })

  afterEach(() => {
    this.sandbox.restore()
  })

  it('should be able to get exam by id', () => {
    const serviceValue = [{ a: 1 }, { a: 2 }]
    const stub = this.sandbox
      .stub(service, 'getMarkedPapersByExamId')
      .resolves(serviceValue)
    const expectedBody = {
      data: [{ a: 1 }, { a: 2 }],
      success: true,
      statusCode: 200
    }
    return request(require('../../../app'))
      .get('/api/v1/exams/id')
      .expect(200)
      .then(res => {
        expect(expectedBody).to.deep.equal(res.body)
        expect(stub.calledWith('id', this.client)).to.be.true
      })
  })

  it('should be able to get paper by id', () => {
    const serviceValue = { a: 1 }
    const stub = this.sandbox
      .stub(service, 'getMarkedPaperByExamIdAndPaperId')
      .resolves(serviceValue)
    const expectedBody = {
      data: { a: 1 },
      success: true,
      statusCode: 200
    }

    return request(require('../../../app'))
      .get('/api/v1/exams/exam-id/papers/paper-id')
      .expect(200)
      .then(res => {
        expect(expectedBody).to.deep.equal(res.body)
        expect(stub.calledWith('exam-id', 'paper-id', this.client)).to.be.true
      })
  })

  it('should be able to send exam', () => {
    return request(require('../../../app'))
      .post('/api/v1/exams')
      .attach('exam', require('path').join(__dirname, './test.json'))
      .expect(201)
      .then(() => {
        expect(this.publishStub.called).to.be.true
      })
  })
})
