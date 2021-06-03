'use strict';

module.exports.info = 'readModel callback';

const helper = require('./helper');
const contractID = 'models';
const version = '1.0';
const name = 'Neural style transfer';

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

    await helper.createModel(bc, ctx, txIndex);

    return Promise.resolve();
}

module.exports.run = function () {
    try {
        const myArgs = {
            chaincodeFunction: 'readModel',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [name + ctx.clientIdx + txIndex.toString()]
        };
        return bc.bcObj.querySmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

module.exports.end = async function () {
    return Promise.resolve();
};
