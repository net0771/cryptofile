import hashlib
import json
import requests
from time import time
from urllib.parse import urlparse


class Blockchain(object):
    def __init__(self):
        self.chain = []
        self.current_transactions = []
        self.nodes = set()

        self.new_block(proof=100, previous_hash=1) # genesis block
    
    def new_block(self, proof, previous_hash=None):
        block = {
            'index': len(self.chain) + 1,
            'timestamp': time(),
            'transactions': self.current_transactions,
            'proof': proof,
            'previous_hash': previous_hash or self.hash(self.chain[-1])
        }

        self.current_transactions = []
        self.chain.append(block)

        return block


    def new_transaction(self, sender, recipient, amount):
        """
        Creteas a transaction -> the next block
        return: index of the next block
        """
        self.current_transactions.append({
            'sender': sender,
            'recipient': recipient,
            'amount': amount
        })

        return self.last_block['index'] + 1


    @staticmethod
    def hash(block):
        """
        creates a SHA256 hash of a Block        
        """
        block_string = json.dumps(block, sort_keys=True).encode()

        return hashlib.sha256(block_string).hexdigest()


    @property
    def last_block(self):
        # returns the last block in the chain(block)
        return self.chain[-1]


    def proof_of_work(self, last_proof):
        proof = 0

        while self.valid_proof(last_proof, proof) is False:
            proof += 1

        return proof


    @staticmethod
    def valid_proof(last_proof, proof):
        guess = str(last_proof * proof).encode()
        guess_hash = hashlib.sha256(guess).hexdigest()

        return guess_hash[:4] == "0000"


    def register_node(self, address):
        """
        add a new node to the list of nodes
        param address: http://127.0.0.1:5000
        """
        
        parse_url = urlparse(address)
        self.nodes.add(parse_url.netloc) # 127.0.0.1:5000


    def valid_chain(self, chain):
        """
        Determine if a given Blockchain is valid
        param chain: a blockchain
        """
        last_block = chain[0]
        current_index = 1
        while current_index < len(chain):
            block = chain[current_index]
            # block check(valid)
            print(f"{last_block}")
            print(f"{block}") # format() 
            print("\n-------------\n")
            if block['previous_hash'] != self.hash(last_block):
                return False

            if not self.valid_proof(last_block['proof'], block['proof']):
                return False

            last_block = block
            current_index += 1

        return True

    def resolve_conflicts(self):
        """
        This is our consensus algorithm
        """
        total_nodes = self.nodes
        new_chain = None

        max_length = len(self.chain)

        for node in total_nodes:
            response = requests.get(f'http://{node}/chain')

            if response.status_code == 200:
                length = response.json()['length']  # 2 
                chain = response.json()['blockchain'] # blocks ...

                if length > max_length and self.valid_chain(chain):
                    max_length = length
                    new_chain = chain
        
        if new_chain:
            self.chain = new_chain
            return True
        
        return False

if __name__ == '__main__':
    bc = Blockchain()
    # index = bc.new_transaction('test_sender', 'test_recipient', 'test_amount')
    # print("Result: " + str(index))
    proof = bc.proof_of_work(100)
    print("PoW=" + str(proof))