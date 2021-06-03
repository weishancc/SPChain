'use strict';

module.exports.info = 'addModel callback';

const contractID = 'models';
const version = '1.0';
const name = 'Neural style transfer';
const address = 'https://hub.docker.com/repository/docker/weishancc/fast_neural_style';
const creator = 'Koma';
const desc = 'Fast neural style transfer';

let bc, ctx;
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

    return Promise.resolve();
};

module.exports.run = function () {
    txIndex++;

    try {
        const myArgs = {
            chaincodeFunction: 'addModel',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [name + ctx.clientIdx + txIndex.toString(), address, creator, desc]
        };
        return bc.bcObj.invokeSmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

module.exports.end = async function () {
    return Promise.resolve();
};
