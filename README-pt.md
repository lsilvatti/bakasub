# BakaSub - O App Completo 🌸✨

[English](README.md)

[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=000000)](https://buymeacoffee.com/leosilvatto)

Hmph! Então você encontrou o repositório de verdade do **BakaSub**, e não só um dos órgãos internos nerdolas dele? Ótimo. Este é o repositório que o usuário final deve realmente rodar na própria máquina. Ele junta frontend, backend e PostgreSQL via Docker Compose para você traduzir legendas sem precisar cuidar manualmente de um servidor Go, um app React e um banco separadamente. B-baka!

> **Presta atenção!** Se você só quer usar o Bakasub, este é o único repositório que importa. Os repositórios de backend e frontend são voltados para desenvolvimento.

## ✨ O Que o Bakasub Faz de Verdade

O Bakasub é um app local para fluxo de legendas focado em tradução assistida por IA e manipulação de trilhas de vídeo.

Com ele você pode:

- inspecionar trilhas de legenda dos seus vídeos
- extrair trilhas de legenda de MKV e outros contêineres suportados
- traduzir arquivos de legenda com modelos do OpenRouter
- enriquecer o contexto da tradução com metadados do TMDB e anotações próprias
- aplicar presets ajustados para anime, filmes, documentários e outros contextos
- reaproveitar traduções em cache para não pagar duas vezes pelas mesmas linhas
- acompanhar jobs, custos, tokens e logs em tempo real
- mesclar a legenda traduzida de volta no seu fluxo de vídeo

## 🌸 Features Que Você Recebe

- **Stack com um comando**: frontend, backend e PostgreSQL sobem juntos com Docker Compose.
- **Interface web moderna**: sem shell Electron, sem drama de instalador desktop separado.
- **Progresso de tradução ao vivo**: a interface acompanha os jobs enquanto o backend trabalha.
- **Memória de tradução**: linhas repetidas podem vir do cache em vez de gastar mais créditos de API.
- **Estimativa pré-voo**: veja lotes, tokens e custo estimado antes de enviar o job.
- **Contexto com TMDB**: associe metadados de filmes ou séries para melhorar a qualidade da tradução.
- **Fluxo baseado em presets**: mantenha estilos separados para anime, filme, comédia e mais.
- **Logs e histórico de jobs**: veja o que aconteceu sem cavar dentro dos containers igual um goblin.

## 🧰 O Que Você Precisa Antes de Começar

Você **não** precisa instalar Go, Node.js ou PostgreSQL no host só para usar o produto.

Você precisa de:

1. **Docker Engine** com o plugin Compose, ou **Docker Desktop**.
2. Uma **release publicada** deste repositório. Prefira uma tag ou GitHub Release, não um commit aleatório de desenvolvimento.
3. Um **caminho absoluto** na sua máquina com os vídeos e legendas que o Bakasub deve enxergar.
4. Uma **chave da OpenRouter API**.
5. Um **token de acesso do TMDB**.

Importante:

- A tradução fica bloqueada até a chave do OpenRouter e o token do TMDB serem configurados e salvos na página de Configurações.
- No setup oficial via Docker deste repositório, FFmpeg e MKVToolNix já vêm dentro da imagem do backend, então você não deve precisar instalar essas ferramentas no host só para começar.

## 🚀 Início Rápido

### 1. Baixe uma release real do Bakasub

Use uma versão taggeada deste repositório.

Por quê? Porque as releases públicas fixam as imagens do backend e frontend em `release.env`. Um checkout em andamento pode ainda conter digests placeholder.

### 2. Crie sua configuração local

Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

Depois edite o `.env` e ajuste pelo menos este valor:

```env
VIDEO_DIR=/caminho/absoluto/para/seus/videos
```

Configurações opcionais que você também pode mudar ali:

- `FRONTEND_HOST_PORT` se a porta `3000` estiver ocupada
- `BACKEND_HOST_PORT` se a porta `8080` estiver ocupada
- `POSTGRES_DB`, `POSTGRES_USER` e `POSTGRES_PASSWORD` se você quiser credenciais locais diferentes para o banco da stack

### 3. Suba a aplicação

```bash
docker compose --env-file release.env --env-file .env up -d
```

Isso sobe:

- o frontend em `http://localhost:3000`
- a API do backend em `http://localhost:8080/api/v1`
- o PostgreSQL interno usado pelo Bakasub

### 4. Abra o Bakasub no navegador

Vá para:

```text
http://localhost:3000
```

### 5. Termine a configuração inicial dentro do app

Abra a página de **Configurações** e configure:

- `OpenRouter API Key`
- `TMDB Access Token`

Salve os dois valores. Até isso acontecer, a página **Traduzir** continua bloqueada.

### 6. Comece a usar o fluxo normal

O fluxo mais comum é:

1. **Extrair**: escolha um vídeo e inspecione ou extraia as trilhas de legenda.
2. **Traduzir**: escolha um arquivo de legenda, um modelo, um preset e um idioma alvo.
3. **Mesclar**: remuxe o resultado traduzido de volta no seu fluxo de vídeo quando precisar.
4. **Logs & Jobs**: revise progresso, erros, custo e histórico.

## 🧭 Uma Boa Primeira Sessão

Se você quiser a primeira execução mais limpa possível, faça isso:

1. Suba a stack.
2. Abra **Configurações** e salve suas credenciais do OpenRouter e do TMDB.
3. Vá em **Modelos e Presets** e marque seus modelos favoritos.
4. Abra **Extrair** e inspecione um vídeo dentro do `VIDEO_DIR` montado.
5. Abra **Traduzir**, associe metadados do TMDB se ajudar, rode a estimativa pré-voo e envie a tradução.
6. Acompanhe o progresso em tempo real e confirme o resultado em **Logs & Jobs**.

## 🗂️ O Que o Bakasub Consegue Enxergar

O Bakasub só enxerga a pasta montada como `VIDEO_DIR`.

Isso significa:

- se um arquivo estiver fora dessa pasta, ele não vai aparecer no app
- se o caminho estiver errado, o Bakasub vai agir como se sua biblioteca estivesse vazia
- se a permissão da pasta estiver bloqueada, o backend não vai conseguir inspecionar nem processar os arquivos

Use caminho absoluto. Caminho relativo não resolve aqui.

## 🔄 Como Atualizar o Bakasub

Quando sair uma nova release pública:

1. Baixe ou atualize para a nova versão taggeada deste repositório.
2. Mantenha seu arquivo `.env` atual.
3. Suba a stack atualizada novamente:

```bash
docker compose --env-file release.env --env-file .env up -d
```

Se quiser ser explícito, pode puxar as imagens antes:

```bash
docker compose --env-file release.env --env-file .env pull
docker compose --env-file release.env --env-file .env up -d
```

## 🛑 Como Parar ou Remover a Stack

Pare os containers sem remover os dados persistentes:

```bash
docker compose --env-file release.env --env-file .env down
```

Remova tudo, incluindo o volume do PostgreSQL usado pelo app:

```bash
docker compose --env-file release.env --env-file .env down -v
```

Só use o segundo comando se você realmente quiser apagar o banco local do app.

## 🩹 Solução de Problemas

### O frontend não abre

- Verifique se a porta `3000` já está em uso.
- Se precisar, altere `FRONTEND_HOST_PORT` no `.env` e suba a stack de novo.

### A porta do backend está ocupada

- Altere `BACKEND_HOST_PORT` no `.env`.
- Reinicie com o mesmo `docker compose --env-file release.env --env-file .env up -d`.

### A tradução está bloqueada

- Abra **Configurações**.
- Salve uma **OpenRouter API Key** válida.
- Salve um **TMDB Access Token** válido.

O produto espera os dois configurados antes de liberar a tradução.

### Não consigo ver meus arquivos de vídeo

- Confirme que `VIDEO_DIR` aponta para o caminho absoluto correto.
- Confirme que o Docker consegue ler essa pasta.
- Confirme que os arquivos realmente estão dentro desse diretório montado.

### A stack não sobe

- Confirme que você está usando uma tag ou release com digests reais em `release.env`.
- Rode `docker compose --env-file release.env --env-file .env logs` para inspecionar a falha.

## 💖 Apoie o Projeto

Se o Bakasub salvou seu cronograma de release, sua fila de legendas ou sua sanidade, você pode apoiar o projeto aqui:

[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=000000)](https://buymeacoffee.com/lsilvatti)