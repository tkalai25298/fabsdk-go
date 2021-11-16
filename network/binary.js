var fs = require('fs'),
    JSON = fs.readFileSync('./channel-artifacts/genesis.block','JSON').toJson(JSON);
process.stdout.write(JSON.substring(0, 500000));