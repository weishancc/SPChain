'use strict';

module.exports.info = 'invokeModel callback';

const contractID = 'models';
const version = '1.0';
const api = 'http://140.123.105.112:5000/mosaic';
const offset = 'https://homepages.cae.wisc.edu/~ece533/images/tulips.png';

let bc, ctx, clientArgs, clientIdx;

/**
* Initializes the workload module before the start of the round.
* @param {BlockchainInterface} blockchain The SUT adapter instance.
* @param {object} context The SUT-specific context for the round.
* @param {object} args The user-provided arguments for the workload module.
*/
module.exports.init = async function (blockchain, context, args) {
    bc = blockchain;
    ctx = context;
    clientArgs = args;
    clientIdx = context.clientIdx.toString();

    return Promise.resolve(); 
}

module.exports.run = function () {
    try {
        const myArgs = {
            chaincodeFunction: 'invokeModel',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [api, offset]
        };
        return bc.bcObj.querySmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

module.exports.end = async function () {
    return Promise.resolve();
};
