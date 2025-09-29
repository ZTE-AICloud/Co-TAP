const _url = '';
const cookie = ``;
const forgerydefense = ``;
const target = _url;

const config = [
  {
    context: ['/api'],
    target: target,
    secure: false,
    ws: true,
    headers: {
      Referer: target,
      Host: target.substring(8),
      forgerydefense: forgerydefense,
      Cookie: cookie,
    },
  },
];
module.exports = config;