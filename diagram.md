```mermaid
---
title: Diagram of GopherCon UK - 2024 Workflow
---
stateDiagram-v2
	direction LR
	
	StatusStarted-->StatusNameCreated
	StatusNameCreated-->StatusColourSet
	StatusColourSet-->StatusAgeDefined
	StatusAgeDefined-->StatusSentToSchool
	StatusAgeDefined-->StatusSentToWork
	StatusSentToSchool-->StatusFinishedSchool
	StatusFinishedSchool-->StatusSentToWork
```