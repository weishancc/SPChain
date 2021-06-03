'use strict';

module.exports.info = 'readLog callback';

const helper = require('./helper');
const contractID = 'logs';
const version = '1.0';
const pk_DS = 'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDFu2UlcTq52DvS+CEzQoXbP+/a\n19Dabb3QGqcPSjpSXrKGnUhS7qiTZadQbiKxu97PWwgRxll0BeHLx/X7lLBXuf50\nGK2/1iv5h0z8EWg2Qx6OqV9OGxmaGaCrMKVgQEK2Df32a/g8Nao5OPWXOv+0Jync\ngYK2WB3wVJS4jsC0hwIDAQAB';
const pk_DC = 'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdn3Pvi6K8oDurzi/eQl7jlhW8\nAj9p+O3QfOxodvuqVCyG7PBdC1zi0+qxTu7sYmGJgcQwtRoyponkmo0lj2wHXQaR\nJMyGPgNslo1Xfrp1bv136ZfickArRed8VTP8v2OL/A/bPRTo29S1uSHRRyEsjoc5\nXVTjpu6IOnKkE9c4kwIDAQAB';
const pk_DP = 'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdn3Pvi6K8oDurzi/eQl7jlhW8\nAj9p+O3QfOxodvuqVCyG7PBdC1zi0+qxTu7sYmGJgcQwtRoyponkmo0lj2wHXQaR\nJMyGPgNslo1Xfrp1bv136ZfickArRed8VTP8v2OL/A/bPRTo29S1uSHRRyEsjoc5\nXVTjpu6IOnKkE9c4kwIDAQAB';

let bc, ctx, limitIndex;
let txIndex = 0;

/**
* Initializes the workload module before the start of the round.
* @param {BlockchainInterface} blockchain The SUT adapter instance.
* @param {object} context The SUT-specific context for the round.
* @param {object} args The user-provided arguments for the workload module.
*/
module.exports.init = async function (blockchain, context, args) {
    bc = blockchain;
    ctx = context;
    limitIndex = args.assets;
    txIndex++;

    await helper.createLog(bc, ctx, txIndex);

    return Promise.resolve();
}

module.exports.run = function () {
    try {
        if (txIndex == limitIndex) {
            txIndex = 0;
        }

        const myArgs = {
            chaincodeFunction: 'readLog',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [pk_DS, pk_DC, pk_DP + ctx.clientIdx + txIndex.toString()]
        };
        return bc.bcObj.querySmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

module.exports.end = async function () {
    return Promise.resolve();
};
