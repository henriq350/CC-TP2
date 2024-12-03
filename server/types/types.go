package types

type Agent struct {
    AgentID string
    AgentIP string
}

// Add agent
func AddAgent(agent Agent, agents map[string]Agent) {
	agents[agent.AgentID] = agent
}

// Remove agent
func RemoveAgent(agent Agent, agents map[string]Agent) {
	delete(agents, agent.AgentID)
}
