# Progetto-DP

Sistema di prenotazione per ristorante sviluppato in Go.

## Struttura del progetto

Il progetto è composto da due moduli Go indipendenti:

- restaurant → servizio principale di prenotazione
- notification → servizio di notifiche email

Ogni modulo ha il proprio `go.mod`.

## Requisiti

- Go >= 1.21
- Git

## Build

Per compilare entrambi i moduli:

```bash
make
