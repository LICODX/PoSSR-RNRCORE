package p2p

// BFT-specific methods for vote and proposal communication

import (
	"encoding/json"
	"fmt"

	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus/bft"
)

// PublishVote broadcasts a BFT vote to the network
func (n *GossipSubNode) PublishVote(vote *bft.Vote) error {
	data, err := json.Marshal(vote)
	if err != nil {
		return fmt.Errorf("failed to marshal vote: %w", err)
	}

	if err := n.voteTopic.Publish(n.ctx, data); err != nil {
		return fmt.Errorf("failed to publish vote: %w", err)
	}

	return nil
}

// PublishProposal broadcasts a BFT proposal to the network
func (n *GossipSubNode) PublishProposal(proposal *bft.Proposal) error {
	data, err := json.Marshal(proposal)
	if err != nil {
		return fmt.Errorf("failed to marshal proposal: %w", err)
	}

	if err := n.proposalTopic.Publish(n.ctx, data); err != nil {
		return fmt.Errorf("failed to publish proposal: %w", err)
	}

	return nil
}

// ListenForVotes listens for incoming BFT votes
func (n *GossipSubNode) ListenForVotes(handler func(*bft.Vote)) {
	go func() {
		for {
			msg, err := n.voteSub.Next(n.ctx)
			if err != nil {
				fmt.Printf("[P2P] Vote subscription error: %v\n", err)
				return
			}

			// Ignore our own messages
			if msg.ReceivedFrom == n.host.ID() {
				continue
			}

			var vote bft.Vote
			if err := json.Unmarshal(msg.Data, &vote); err != nil {
				fmt.Printf("[P2P] Failed to decode vote: %v\n", err)
				continue
			}

			// Call handler
			handler(&vote)
		}
	}()
}

// ListenForProposals listens for incoming BFT proposals
func (n *GossipSubNode) ListenForProposals(handler func(*bft.Proposal)) {
	go func() {
		for {
			msg, err := n.proposalSub.Next(n.ctx)
			if err != nil {
				fmt.Printf("[P2P] Proposal subscription error: %v\n", err)
				return
			}

			// Ignore our own messages
			if msg.ReceivedFrom == n.host.ID() {
				continue
			}

			var proposal bft.Proposal
			if err := json.Unmarshal(msg.Data, &proposal); err != nil {
				fmt.Printf("[P2P] Failed to decode proposal: %v\n", err)
				continue
			}

			// Call handler
			handler(&proposal)
		}
	}()
}
