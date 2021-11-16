const fs = require('fs')
fs.readFile('./string.txt' , (err, data) => {
  if (err) {
    console.error(err)
    return }
  const obj = {
      key : data.toString()
  }
  console.log(JSON.stringify(obj))
})