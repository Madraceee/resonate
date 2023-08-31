package network

import (
	"math/rand"
	"net/http"

	"github.com/resonatehq/resonate/internal/aio"
	"github.com/resonatehq/resonate/internal/kernel/bus"
	"github.com/resonatehq/resonate/internal/kernel/types"
	"github.com/resonatehq/resonate/internal/util"
)

type NetworkDST struct {
	r *rand.Rand
	p float32
}

type NetworkDSTDevice struct {
	*NetworkDST
}

func NewDST(r *rand.Rand, p float32) aio.Subsystem {
	return &NetworkDST{
		r: r,
		p: p,
	}
}

func (n *NetworkDST) String() string {
	return "network:dst"
}

func (n *NetworkDST) Start() error {
	return nil
}

func (n *NetworkDST) Stop() error {
	return nil
}

func (n *NetworkDST) Reset() error {
	return nil
}

func (n *NetworkDST) NewWorker(int) aio.Worker {
	return &NetworkDSTDevice{n}
}

func (d *NetworkDSTDevice) Process(sqes []*bus.SQE[types.Submission, types.Completion]) []*bus.CQE[types.Submission, types.Completion] {
	cqes := make([]*bus.CQE[types.Submission, types.Completion], len(sqes))

	for i, sqe := range sqes {
		util.Assert(sqe.Submission.Network != nil, "submission must not be nil")

		switch sqe.Submission.Network.Kind {
		case types.Http:
			cqe := &bus.CQE[types.Submission, types.Completion]{
				Kind:     sqe.Kind,
				Callback: sqe.Callback,
			}

			var res *http.Response

			if d.r.Float32() < d.p {
				res = &http.Response{
					StatusCode: http.StatusOK,
				}
			} else {
				res = &http.Response{
					StatusCode: http.StatusInternalServerError,
				}
			}

			cqe.Completion = &types.Completion{
				Kind: types.Network,
				Network: &types.NetworkCompletion{
					Kind: types.Http,
					Http: res,
				},
			}

			cqes[i] = cqe
		default:
			panic("invalid network submission")
		}
	}

	return cqes
}
