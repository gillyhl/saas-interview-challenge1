const fs = require('fs')
const path = require('path')

const POSSIBLE_ANSWERS = 'ABCDE'

const generateRandomAnswer = () => POSSIBLE_ANSWERS[Math.floor(Math.random() * POSSIBLE_ANSWERS.length)]

const generate = (name, numberOfQuestions, numberOfPapers, probability, fileName) => {
  const numberOfQuestionsInt = parseInt(numberOfQuestions)
  const numberOfPapersInt = parseInt(numberOfPapers)
  const probabilityInt = parseInt(probability)
  const answers = Array.apply(null, {length: numberOfQuestionsInt})
  .map((_, i) => ({
    questionNumber: i + 1,
    answer: generateRandomAnswer()
  }))
  const exam = {
    name,
    answers,
    papers: Array.apply(null, {length: numberOfPapersInt})
    .map((_, i) => ({
      id: `student-${i}`,
      answers: answers.map(answer => ({
        ...answer,
        ...(Math.random() * 100) > probabilityInt ? { answer: generateRandomAnswer() } : {}
      }))
    }))
    
  }

  const jsonContent = JSON.stringify(exam, null, 2)
  fs.writeFile(path.join(__dirname, fileName), jsonContent, 'utf8', function (err) {
    if (err) {
        console.log("An error occured while writing JSON Object to File.");
        return console.log(err);
    }
 
    console.log("JSON file has been saved.");
  })
}

generate(...process.argv.slice(2))