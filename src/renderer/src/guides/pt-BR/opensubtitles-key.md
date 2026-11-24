# Como conseguir uma API key do OpenSubtitles

Quando um vídeo não tem faixa de legenda para traduzir, o BakaSub pode buscar uma no [OpenSubtitles](https://www.opensubtitles.com/) e traduzir a partir dela. Isso é opcional: você também pode escolher um arquivo de legenda que já tem.

Cada pessoa usa a própria API key, e ela é gratuita.

## Passos

1. Crie uma conta em [opensubtitles.com](https://www.opensubtitles.com/) e confirme o e-mail.
2. Abra a página **API consumers** em [opensubtitles.com/consumers](https://www.opensubtitles.com/consumers).
3. Crie um consumer novo. Basta dar um nome, por exemplo `BakaSub`.
4. Copie a API key que aparece para esse consumer.
5. No BakaSub, cole em **Configurações → Chaves → OpenSubtitles**. O BakaSub confere a key com o OpenSubtitles antes de salvar.

O OpenSubtitles pode mudar o nome de botões e páginas com o tempo. A página de consumers fica sempre dentro da sua conta.

## Limite de downloads e a sua conta

O OpenSubtitles limita quantas legendas você pode baixar por dia. Buscar não conta, só baixar. O BakaSub guarda o que baixou, então escolher a mesma legenda de novo não gasta nada.

Se quiser um limite maior, entre com a sua conta do OpenSubtitles na mesma seção de Configurações. O limite passa a ser o da conta. Use o nome de usuário, porque o OpenSubtitles não aceita o e-mail no login.

## Qual resultado escolher

- Prefira os resultados marcados com "Feita para este arquivo". Essa legenda foi feita para este arquivo de vídeo exato, então o tempo bate.
- Os outros resultados batem só pelo título e podem estar fora de sincronia com o seu vídeo. O BakaSub não ajusta o tempo.
- Evite os marcados com "Tradução automática". Traduzir uma tradução perde sentido duas vezes.

## Privacidade

- A API key e a senha ficam criptografadas no seu computador, no chaveiro do sistema quando ele existe. Elas só saem dele nas requisições ao OpenSubtitles.
- Uma busca envia o título, a temporada e o episódio tirados do nome do arquivo. Envia também uma impressão digital do arquivo de vídeo, calculada a partir do tamanho e de um pequeno trecho do conteúdo. O vídeo em si nunca é enviado.
