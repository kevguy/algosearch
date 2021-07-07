const constants = require('./service/global1');
const algosdk = require('algosdk');

const algoUrl = new URL(constants.algodurl);
const client = new algosdk.Algodv2(
    constants.algodapi,
    algoUrl,
    algoUrl.port ? algoUrl.port : 8080);

function Int64ToString(bytes, isSigned) {
    const isNegative = isSigned && bytes.length > 0 && bytes[0] >= 0x80;
    const digits = [];
    bytes.forEach((byte, j) => {
        if(isNegative)
            byte = 0x100 - (j == bytes.length - 1 ? 0 : 1) - byte;
        for(let i = 0; byte > 0 || i < digits.length; i++) {
            byte += (digits[i] || 0) * 0x100;
            digits[i] = byte % 10;
            byte = (byte - digits[i]) / 10;
        }
    });
    return (isNegative ? '-' : '') + digits.reverse().join('');
}

const tests = [
    {
        inp: [77, 101, 130, 33, 7, 252, 253, 82],
        signed: false,
        expectation: '5577006791947779410'
    },
    {
        inp: [255, 255, 255, 255, 255, 255, 255, 255],
        signed: true,
        expectation: '-1'
    },
];

tests.forEach(test => {
    const result = Int64ToString(test.inp, test.signed);
    console.log(`${result} ${result !== test.expectation ? '!' : ''}=== ${test.expectation}`);
});

async function test(blockNum) {
    const blk = await client.block(blockNum).do();
    const proposer = algosdk.encodeAddress(blk["cert"]["prop"]["oprop"]);

    console.log('proposer');
    console.log(proposer);

    console.log('blockHash');
    console.log(blk["cert"]["prop"]["dig"]);

    // const blockHash = algosdk.encodeObj(blk["cert"]["prop"]["dig"]);
    try {
        console.log('Using algosdk.encodeUint64');
        const blockHash = algosdk.encodeUint64(blk["cert"]["prop"]["dig"]);
        console.log('blockHash');
        console.log(blockHash);
    } catch (e) {
        console.log('Failed');
        console.log(e);
    }

    try {
        console.log('Using algosdk.encodeAddress');
        const blockHash = algosdk.encodeAddress(blk["cert"]["prop"]["dig"]);
        console.log('blockHash');
        console.log(blockHash);
    } catch (e) {
        console.log('Failed');
        console.log(e);
    }

    try {
        console.log('Using algosdk.encodeObj');
        const blockHash = algosdk.encodeObj(blk["cert"]["prop"]["dig"]);
        console.log('blockHash');
        console.log(blockHash);
    } catch (e) {
        console.log('Failed');
        console.log(e);
    }

    try {
        console.log('Using algosdk.encodeUnsignedTransaction');
        const blockHash = algosdk.encodeUnsignedTransaction(blk["cert"]["prop"]["dig"]);
        console.log('blockHash');
        console.log(blockHash);
    } catch (e) {
        console.log('Failed');
        console.log(e);
    }

    try {
        console.log('using custom function');
        const blockHash = Int64ToString(blk["cert"]["prop"]["dig"], true);
        console.log('blockhash');
        console.log(blockHash);
    } catch (e) {
        console.log('Failed');
        console.log(e);
    }

    try {
        console.log('using custom function');
        const blockHash = Int64ToString(blk["cert"]["prop"]["dig"], false);
        console.log('blockhash');
        console.log(blockHash);
    } catch (e) {
        console.log('Failed');
        console.log(e);
    }

    try {
        console.log('base64');
        // const blockHash = Buffer.from(blk["cert"]["prop"]["dig"], 'base64').toString();
        const blockHash = Buffer.from(blk["cert"]["prop"]["dig"]).toString("base64");
        console.log('blockhash');
        console.log(blockHash);
    } catch (e) {
        console.log('Failed');
        console.log(e);
    }
}


test(11)
    .then(() => { return client.block(12).do(); })
    .then((res) => { console.log(res); });
