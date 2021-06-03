'use strict';

module.exports.info = 'addWallet callback';

const contractID = 'wallets';
const version = '1.0';
const name = 'Koma';
const price = '1000000';

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
}

module.exports.run = function () {
    txIndex++;

    try {
        const myArgs = {
            chaincodeFunction: 'addWallet',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [name + ctx.clientIdx + txIndex.toString(), price]
        };
        return bc.bcObj.invokeSmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

module.exports.end = async function () {
    return Promise.resolve();
};
