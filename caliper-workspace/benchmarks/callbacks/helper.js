'use strict';

/**
* artworks helper
*/
var contractID = 'artworks';
const version = '1.0';
const tokenID = 'afee1720-28ef-4790-8fe8-50f38681b919';
const multiHash = 'zdpuAqazmCjGhntRXpsKpcsUou9D72TzowCTgJK5jC62CoYvt';
const owner = 'Koma';
const creator = 'Koma';

module.exports.createArtwork = async function (bc, ctx, txIndex) {
    try {
        const myArgs = {
            chaincodeFunction: 'uploadArtwork',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [tokenID + ctx.clientIdx + txIndex.toString(), multiHash, owner, creator]
        };
        await bc.bcObj.invokeSmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

/**
* logs helper
*/
contractID = 'logs'
const pk_DS = 'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDFu2UlcTq52DvS+CEzQoXbP+/a\n19Dabb3QGqcPSjpSXrKGnUhS7qiTZadQbiKxu97PWwgRxll0BeHLx/X7lLBXuf50\nGK2/1iv5h0z8EWg2Qx6OqV9OGxmaGaCrMKVgQEK2Df32a/g8Nao5OPWXOv+0Jync\ngYK2WB3wVJS4jsC0hwIDAQAB';
const pk_DC = 'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdn3Pvi6K8oDurzi/eQl7jlhW8\nAj9p+O3QfOxodvuqVCyG7PBdC1zi0+qxTu7sYmGJgcQwtRoyponkmo0lj2wHXQaR\nJMyGPgNslo1Xfrp1bv136ZfickArRed8VTP8v2OL/A/bPRTo29S1uSHRRyEsjoc5\nXVTjpu6IOnKkE9c4kwIDAQAB';
const pk_DP = 'MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdn3Pvi6K8oDurzi/eQl7jlhW8\nAj9p+O3QfOxodvuqVCyG7PBdC1zi0+qxTu7sYmGJgcQwtRoyponkmo0lj2wHXQaR\nJMyGPgNslo1Xfrp1bv136ZfickArRed8VTP8v2OL/A/bPRTo29S1uSHRRyEsjoc5\nXVTjpu6IOnKkE9c4kwIDAQAB';
const sk_data = 'MIICXQIBAAKBgQDFu2UlcTq52DvS+CEzQoXbP+/a19Dabb3QGqcPSjpSXrKGnUhS\n7qiTZadQbiKxu97PWwgRxll0BeHLx/X7lLBXuf50GK2/1iv5h0z8EWg2Qx6OqV9O\nGxmaGaCrMKVgQEK2Df32a/g8Nao5OPWXOv+0JyncgYK2WB3wVJS4jsC0hwIDAQAB\nAoGAGoUU6U3Y9bER66yaQ1ZCPr2XYronCbdwnwHcHdmzncrpVdMFhBNR/8bsBwZm\ngoEIlAteYd1BBSXUJXE8cCbT9GDzZlTiXRcHOuFanq382cgmQKnBkRv60VH0zICj\nkE7fWLE9wpHHKuApsOhFSuom63aXOaaqvVee+DROi1R7FBECQQDPvJOqbB7gttxb\nMC2yTAi6pP1kZHBYWJD6rQmb34M0LonVc/GzsApebCpYSsx8n2ZG873mB9kPpDyP\nWqK430UxAkEA86vHqCZ5x5xn2GEhoXH5ztHDF3OZ5jDorEtncgt3xLWc9Z87QmN8\nGohjBQI/Y51ErBKteeSOoGhtPzjwywyHNwJBAMcCoXRioDIm/HNfdGea78HezeGf\nVwFL15hOrSXmuosDCoiyypqZy1UpymdLQRsimZjfaM02N3wEmv+6lKkHPAECQQCY\nrcHkcndLw4yt3+6aojfMh1KelyiPO4YOrxCaPOVGtCUtIiCXcI6KcXrZ4JanbBtj\nVjCsd7GGgOgy/RKjp63xAkB6m7dlRCat1R8HYCorESAtmjGLzD9B2gKnczhYj+Uh\nnea68vDSW8HXZVdxmY8X3SGJB6Og9XQu+Jg/Sl7JnU5I';

module.exports.createLog = async function (bc, ctx, txIndex) {
    try {
        const myArgs = {
            chaincodeFunction: 'addLog',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [pk_DS, pk_DC, pk_DP + ctx.clientIdx + txIndex.toString(), sk_data, '1', 'R']
        };
        await bc.bcObj.invokeSmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

/**
* models helper
*/
contractID = 'models'
const name = 'Neural style transfer';
const address = 'https://hub.docker.com/repository/docker/weishancc/fast_neural_style';
const desc = 'Fast neural style transfer';

module.exports.createModel = async function (bc, ctx, txIndex) {
    try {
        const myArgs = {
            chaincodeFunction: 'addModel',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [name + ctx.clientIdx + txIndex.toString(), address, creator, desc]
        };
        await bc.bcObj.invokeSmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

/**
* wallets helper
*/
contractID = 'wallets';
const balance = '1000000';

module.exports.createWallet = async function (bc, ctx, txIndex, role) {
    try {
        const myArgs = {
            chaincodeFunction: 'addWallet',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [role + ctx.clientIdx + txIndex.toString(), balance]
        };

        await bc.bcObj.invokeSmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${clientIdx}: Smart Contract threw with error: ${error}`);
    };
};