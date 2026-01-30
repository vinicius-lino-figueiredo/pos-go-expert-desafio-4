**Objetivo**: Desenvolvedor  um sistema Go que receba um CEP, idenfica a cidade
e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin).
Esse sistema deverá ser publicado no Google Cloud Run.
**Requisitos:**
- O sistema deve receber um CEP válido de 8 dígitos
- O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização,
  a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius
  Fahrenheit, Kelvin.
- O sistema deve responder adequadamente nos seguintes cenários:
  - Em caso de sucesso:
    - Código HTTP:**200**
    - Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
  - Em caso de falha, caso o CEP não seja válido (com formato correto):
    - Código HTTP: 422
    - Mensagem: **invalid zipcode**
  - Em caso de falha, caso o CEP não seja encontrado:
   - Código HTTP: **404**
   - Mensagem: **can not find zipcode**
- Deverá ser realizado o deploy no Google Cloud Run.
**Dicas:**
- Utilize a API viaCEP (ou similar) para encontrar a localização que deseja
  consultar a temperatura: <https://viacep.com.br/>
- Utiliza a API WheatherAPI (ou similar) para consular as temperaturas
  desejadas: <https://www.weatherapi.com/>
Para realizar a conversão de Celsius para Fahrenheit, utilize aseguinte fórmula:
  `F = C * 1,8 + 32`
- Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula:
  `K = C + 273`
  - Sendo F = Fahrenheit
  - Sendo C = Celsius
  - Sendo K = Kelvin
**Entrega:**
- O código-fonte completo da implementação
- Testes automatizados demonstrando o funcionamento
- Utilize docker/docker-compose para que possamos realizar os testes de sua
 aplicação.
- Deploy realizado no Google Cloud Run (free tier) e o endereço ativo para ser
  acessado.

---

**URL da aplicação publicada:**

https://desafio-4-763797664904.southamerica-east1.run.app/temperature/{cep}
