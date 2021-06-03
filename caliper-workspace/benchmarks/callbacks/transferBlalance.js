'use strict';

module.exports.info = 'transferBalance callback';

const contractID = 'wallets';
const version = '1.0';
const payer = 'Amber';
const payee = 'Koma';
const price = '100';
const creator = 'Koma';
const r_y = '0.05';

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

    return Promise.resolve();
}

module.exports.run = function () {
    try {
         const myArgs = {
            chaincodeFunction: 'transferBalance',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [payer + ctx.clientIdx + txIndex.toString(), payee + ctx.clientIdx + txIndex.toString(), creator + ctx.clientIdx + txIndex.toString(), price, r_y]
        };
        return bc.bcObj.invokeSmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

module.exports.end = async function () {
    return Promise.resolve();
};
