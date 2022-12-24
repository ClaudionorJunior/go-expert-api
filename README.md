# Descrição do projeto:
Esse projeto faz parte do curso Go Expert da [FullCycle](https://fullcycle.com.br/). Foi desevolvido a partir dos aprendizados de como construir uma API com Web Framework, bem como testes e lidar com pacotes.<br>
É uma API que possui um CRUD de `products`, criação de usuário e autenticação de usuário com `acessToken`.
<br><br>
<img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/ClaudionorJunior/go-expert-api">

# Rodando o projeto:
```sh
cd cmd/server && go run main.go
```

# Atualizando docs:
```sh
swag init -g cmd/server/main.go -o api/
```

# Acessando docs(Swagger):

```sh
http://localhost:8000/docs/index.html
```

## Autor
<table>
  <tr>
    <th><img src="https://avatars.githubusercontent.com/u/82416762?v=4" width=60></th>
  </tr>
  <tr>
    <td><a href="https://github.com/ClaudionorJunior">Github</a></td>
  </tr>
  <tr>
    <td><a href="https://www.linkedin.com/in/claudionorsilva">Linkedin</a></td>
  </tr>
</table>
