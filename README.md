# Open Telemetry

- É Um framework de observabilidade agnostico a vendors, para softwares cloud native. Ele oferece o necessário para trabalharmos com observabilidade.

- É um conjunto de ferramentas que disponibiliza APIs e SDKs

- Ele trabalha com instrumentação, geração, coleta e exportação de dados de telemetria

## Collector

- É um agente ou serviço, que pega os dados de telemetria gerados para serem enviados para os vendors;
- É um pipeline de processamento para envio dos dados para os vendors;
- É agnostico ao vendor;

## Instrumentação

- É uma forma de lidarmos com os comportamentos de uma aplicação, para a partir disso conseguirmos gerar dados para serem enviados para um collector do open telemetry ou diretamente para um vendor

Observabilidade:
    - logs
    - metricas
    - tracing