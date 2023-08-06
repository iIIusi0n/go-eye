package main

var (
	ships = []string{
		"Rifter",
		"Atron",
		"Flycatcher",
		"Stiletto",
		"Vargur",
		"Muninn",
		"Condor",
	}

	fakeMsg = "Loss Count: 7\n" +
		"<Prop Analysis>\n" +
		"AB/MWD/DUAL PROP/NO PROP: 5/1/0/2\n" +
		"Prop uses: MWD AB AB AB AB AB NP NP\n" +
		"\n" +
		"<Weapon Analysis>\n" +
		"Long Range/Short Range: 2/6\n" +
		"Weapon uses: LR LR SR SR SR SR SR SR\n" +
		"\n" +
		"<Web/Scram Analysis>\n" +
		"Web/Scram: 8/8\n" +
		"EWAR uses: WS WS WS WS WS WS WWS WWS\n"
)
