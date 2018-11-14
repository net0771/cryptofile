web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:8545"));
abi = JSON.parse('[{"constant":false,"inputs":[{"name":"candidate","type":"bytes32"}],"name":"totalVotesFor","outputs":[{"name":"","type":"uint8"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"candidate","type":"bytes32"}],"name":"validCandidate","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"votesReceived","outputs":[{"name":"","type":"uint8"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"x","type":"bytes32"}],"name":"bytes32ToString","outputs":[{"name":"","type":"string"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"candidateList","outputs":[{"name":"","type":"bytes32"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"candidate","type":"bytes32"}],"name":"voteForCandidate","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"contractOwner","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"inputs":[{"name":"candidateNames","type":"bytes32[]"}],"payable":false,"type":"constructor"}]');
VotingContext = web3.eth.contract(abi);
contractInstance = VotingContext.at('0x2d464e9fcb10d2a03c5768c379e69ea51a7e809a');
candidates = {"0xb7778275F81fC85D4030752aE9dfbbaD66629C6d": "candidate-1",
              "0xfBC666302d9D64be97642e2cbD8b8c94d45b61a6": "candidate-2",
              "0xB8f3018A53e8CC8eb88bf24d588e6609C2ae9794": "candidate-3"};

function voteForCandidate() {
  candidateName = $("#candidate").val();
  
  contractInstance.voteForCandidate(candidateName, {gas: 140000, from: web3.eth.accounts[0]}, function() {
    let div_id = candidates[candidateName];
    let toalVote = contractInstance.totalVotesFor.call(candidateName).toString();
    $("#" + div_id).html(toalVote);
    console.log("#" + candidateName +" has been voted [" + toalVote + "].");
  });
}

$(document).ready(function(){
    candidateNames = Object.keys(candidates);

    for (var i = 0; i < candidateNames.length; i++) {
        let name = candidateNames[i];
        console.log(name);
        let val = contractInstance.totalVotesFor.call(name).toString();
        $('#' + candidates[name]).html(val);
    }
});