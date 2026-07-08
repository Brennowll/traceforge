package simulation

func (r *SimulationResult) addTrace(event TraceEvent) {
	r.Trace = append(r.Trace, event)
}
