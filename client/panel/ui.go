package panel

// panelHTML is the complete single-page management UI.
// It is served at /panel/ and communicates with /panel/api/*.
const panelHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>frpc panel</title>
<style>
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;500;600&family=Syne:wght@400;500;600;700;800&display=swap');
:root{
  --bg:#07090f;--sur:#0d1018;--sur2:#131720;--bdr:#181d28;--bdr2:#1e2535;
  --acc:#2de8b0;--acc2:#4dabff;--acc3:#ff5f6d;--acc4:#ffcc44;
  --txt:#d8e8ff;--txt2:#6b7a99;--txt3:#303a52;
  --mono:'JetBrains Mono',monospace;--sans:'Syne',sans-serif;
}
*{margin:0;padding:0;box-sizing:border-box;}
body{font-family:var(--sans);background:var(--bg);color:var(--txt);min-height:100vh;}
body::before{content:'';position:fixed;inset:0;pointer-events:none;z-index:0;
  background:radial-gradient(ellipse 70% 50% at 15% -10%,rgba(45,232,176,.07),transparent),
             radial-gradient(ellipse 50% 40% at 85% 105%,rgba(77,171,255,.05),transparent);}

/* LOGIN */
#login{position:fixed;inset:0;display:flex;align-items:center;justify-content:center;z-index:100;background:var(--bg);}
#app{display:none;min-height:100vh;position:relative;z-index:1;}
.lw{width:360px;}
.lb{display:flex;align-items:center;gap:12px;margin-bottom:44px;}
.li{width:44px;height:44px;background:linear-gradient(135deg,var(--acc),var(--acc2));border-radius:10px;
  display:flex;align-items:center;justify-content:center;font-size:18px;font-weight:800;color:#000;font-family:var(--mono);}
.ln{font-size:22px;font-weight:800;letter-spacing:-.5px;}
.ln em{color:var(--acc);font-style:normal;}
.ls{font-size:10px;color:var(--txt2);font-family:var(--mono);letter-spacing:2px;margin-bottom:32px;}
.fl{margin-bottom:16px;}
.fl label{display:block;font-size:10px;font-weight:600;color:var(--txt2);text-transform:uppercase;letter-spacing:1.5px;margin-bottom:7px;font-family:var(--mono);}
.fi,.fs,.fta{width:100%;background:var(--sur);border:1px solid var(--bdr2);border-radius:7px;
  padding:10px 13px;color:var(--txt);font-family:var(--mono);font-size:13px;outline:none;transition:border-color .2s;}
.fi:focus,.fs:focus,.fta:focus{border-color:var(--acc);}
.fi::placeholder,.fta::placeholder{color:var(--txt3);}
.fs{appearance:none;cursor:pointer;padding-right:30px;
  background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='6'%3E%3Cpath d='M1 1l4 4 4-4' stroke='%23303a52' stroke-width='1.5' fill='none'/%3E%3C/svg%3E");
  background-repeat:no-repeat;background-position:right 11px center;}
.fs option{background:var(--sur);}
.fta{resize:vertical;min-height:160px;line-height:1.7;font-size:12px;}
.btn{display:inline-flex;align-items:center;justify-content:center;gap:6px;
  padding:10px 20px;border:none;border-radius:7px;font-family:var(--sans);
  font-size:13px;font-weight:700;cursor:pointer;transition:all .15s;letter-spacing:.3px;}
.bp{background:var(--acc);color:#000;width:100%;}
.bp:hover{filter:brightness(1.1);transform:translateY(-1px);}
.bg{background:transparent;border:1px solid var(--bdr2);color:var(--txt2);}
.bg:hover{border-color:var(--acc);color:var(--acc);}
.bd{background:rgba(255,95,109,.12);border:1px solid rgba(255,95,109,.25);color:var(--acc3);}
.bd:hover{background:rgba(255,95,109,.2);}
.bsm{padding:7px 13px;font-size:11px;border-radius:6px;}
.bic{background:transparent;border:1px solid var(--bdr2);color:var(--txt2);
  padding:6px;border-radius:6px;cursor:pointer;transition:all .15s;
  font-size:14px;display:inline-flex;align-items:center;line-height:1;}
.bic:hover{border-color:var(--acc);color:var(--acc);}
.emsg{background:rgba(255,95,109,.1);border:1px solid rgba(255,95,109,.2);
  border-radius:6px;padding:9px 13px;font-size:12px;color:var(--acc3);margin-top:12px;display:none;font-family:var(--mono);}

/* LAYOUT */
.side{position:fixed;left:0;top:0;bottom:0;width:222px;background:var(--sur);
  border-right:1px solid var(--bdr);display:flex;flex-direction:column;z-index:50;}
.sh{padding:22px 18px;border-bottom:1px solid var(--bdr);display:flex;align-items:center;gap:10px;}
.si{width:34px;height:34px;background:linear-gradient(135deg,var(--acc),var(--acc2));
  border-radius:7px;display:flex;align-items:center;justify-content:center;
  font-size:13px;font-weight:800;color:#000;font-family:var(--mono);}
.st{font-size:16px;font-weight:800;letter-spacing:-.3px;}
.st em{color:var(--acc);font-style:normal;}
.snav{flex:1;padding:14px 10px;display:flex;flex-direction:column;gap:2px;}
.ni{display:flex;align-items:center;gap:9px;padding:9px 11px;border-radius:7px;
  cursor:pointer;color:var(--txt2);font-size:13px;font-weight:600;transition:all .15s;
  border:none;background:none;width:100%;text-align:left;}
.ni:hover{background:var(--sur2);color:var(--txt);}
.ni.on{background:rgba(45,232,176,.1);color:var(--acc);}
.nic{font-size:15px;width:20px;text-align:center;}
.sf{padding:12px 10px;border-top:1px solid var(--bdr);}
.sp{display:flex;align-items:center;gap:8px;padding:8px 11px;border-radius:6px;
  background:var(--sur2);margin-bottom:7px;font-size:11px;font-family:var(--mono);}
.dot{width:7px;height:7px;border-radius:50%;background:var(--txt3);flex-shrink:0;}
.dot.r{background:var(--acc);box-shadow:0 0 6px var(--acc);animation:pu 2s infinite;}
.dot.e{background:var(--acc3);}
.dot.w{background:var(--acc4);}
@keyframes pu{0%,100%{opacity:1}50%{opacity:.3}}
.main{margin-left:222px;min-height:100vh;display:flex;flex-direction:column;}
.top{padding:15px 28px;border-bottom:1px solid var(--bdr);display:flex;
  align-items:center;justify-content:space-between;background:var(--sur);position:sticky;top:0;z-index:30;}
.pt{font-size:17px;font-weight:800;letter-spacing:-.3px;}
.ps{font-size:10px;color:var(--txt2);font-family:var(--mono);margin-top:1px;letter-spacing:.5px;}
.con{padding:24px 28px;flex:1;}
.sec{display:none;}.sec.on{display:block;}

/* STATS */
.sg{display:grid;grid-template-columns:repeat(4,1fr);gap:12px;margin-bottom:18px;}
.sc{background:var(--sur);border:1px solid var(--bdr);border-radius:9px;padding:18px;position:relative;overflow:hidden;}
.sc::after{content:'';position:absolute;bottom:0;left:0;right:0;height:2px;}
.sc.c0::after{background:var(--acc2)}.sc.c1::after{background:var(--acc)}.sc.c2::after{background:var(--acc3)}.sc.c3::after{background:var(--acc4)}
.sl{font-size:10px;color:var(--txt2);text-transform:uppercase;letter-spacing:1.5px;font-family:var(--mono);margin-bottom:7px;}
.sv{font-size:30px;font-weight:800;font-family:var(--mono);}
.sc.c0 .sv{color:var(--acc2)}.sc.c1 .sv{color:var(--acc)}.sc.c2 .sv{color:var(--acc3)}.sc.c3 .sv{color:var(--acc4)}

/* CARD */
.card{background:var(--sur);border:1px solid var(--bdr);border-radius:9px;overflow:hidden;margin-bottom:18px;}
.ch{padding:14px 20px;border-bottom:1px solid var(--bdr);display:flex;align-items:center;justify-content:space-between;}
.ct{font-size:13px;font-weight:700;display:flex;align-items:center;gap:7px;}
.cb{padding:20px;}

/* CTRL */
.cb2{display:flex;gap:9px;flex-wrap:wrap;margin-bottom:18px;}
.cc{display:flex;align-items:center;gap:7px;padding:9px 16px;border-radius:7px;
  border:1px solid var(--bdr2);background:var(--sur2);color:var(--txt);
  font-family:var(--sans);font-size:12px;font-weight:600;cursor:pointer;transition:all .15s;}
.cc:hover{border-color:var(--acc2);color:var(--acc2);}
.cc.cs:hover{border-color:var(--acc);color:var(--acc);}
.cc.cx:hover{border-color:var(--acc3);color:var(--acc3);}
.cc.cr:hover{border-color:var(--acc4);color:var(--acc4);}

/* TABLE */
.tbl{width:100%;border-collapse:collapse;font-size:12px;}
.tbl th{padding:10px 14px;text-align:left;font-size:10px;font-weight:600;color:var(--txt3);
  text-transform:uppercase;letter-spacing:1.5px;font-family:var(--mono);border-bottom:1px solid var(--bdr);}
.tbl td{padding:12px 14px;border-bottom:1px solid var(--bdr);font-family:var(--mono);font-size:12px;}
.tbl tr:last-child td{border-bottom:none;}
.tbl tr:hover td{background:var(--sur2);}
.tbl tr.dis td{opacity:.35;}
.pn{font-weight:600;color:var(--txt);}
.ag{display:flex;gap:5px;}
.bdg{display:inline-flex;align-items:center;padding:2px 7px;border-radius:4px;
  font-size:10px;font-weight:600;font-family:var(--mono);text-transform:uppercase;letter-spacing:.5px;}
.br{background:rgba(45,232,176,.12);color:var(--acc);border:1px solid rgba(45,232,176,.18);}
.be{background:rgba(255,95,109,.12);color:var(--acc3);border:1px solid rgba(255,95,109,.18);}
.bd2{background:rgba(48,58,82,.3);color:var(--txt3);border:1px solid var(--bdr2);}
.bt{background:rgba(77,171,255,.12);color:var(--acc2);border:1px solid rgba(77,171,255,.18);}
.bu{background:rgba(255,204,68,.12);color:var(--acc4);border:1px solid rgba(255,204,68,.18);}
.bh{background:rgba(45,232,176,.08);color:var(--acc);border:1px solid rgba(45,232,176,.12);}

/* TOGGLE */
.tog{position:relative;width:36px;height:20px;display:inline-block;}
.tog input{opacity:0;width:0;height:0;}
.tsl{position:absolute;inset:0;background:var(--bdr2);border-radius:20px;cursor:pointer;transition:.2s;}
.tsl::before{content:'';position:absolute;width:14px;height:14px;left:3px;top:3px;
  background:var(--txt2);border-radius:50%;transition:.2s;}
.tog input:checked+.tsl{background:rgba(45,232,176,.22);}
.tog input:checked+.tsl::before{transform:translateX(16px);background:var(--acc);}

/* GRID */
.g2{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:12px;}
.g1{margin-bottom:12px;}

/* MODAL */
.ov{position:fixed;inset:0;background:rgba(0,0,0,.75);backdrop-filter:blur(4px);
  z-index:200;display:none;align-items:center;justify-content:center;}
.ov.on{display:flex;}
.modal{background:var(--sur);border:1px solid var(--bdr2);border-radius:11px;
  width:520px;max-height:90vh;overflow-y:auto;animation:su .2s ease;}
@keyframes su{from{opacity:0;transform:translateY(14px)}to{opacity:1;transform:none}}
.mh{padding:20px 22px;border-bottom:1px solid var(--bdr);display:flex;
  align-items:center;justify-content:space-between;position:sticky;top:0;background:var(--sur);}
.mtl{font-size:15px;font-weight:700;}
.mcl{background:none;border:none;color:var(--txt2);cursor:pointer;font-size:18px;transition:color .15s;}
.mcl:hover{color:var(--txt);}
.mb{padding:20px 22px;}
.mft{padding:13px 22px;border-top:1px solid var(--bdr);display:flex;gap:7px;justify-content:flex-end;}

/* NOTIFY */
.notif{position:fixed;top:16px;right:16px;background:var(--sur);
  border:1px solid var(--bdr2);border-radius:8px;padding:12px 16px;min-width:220px;
  z-index:9999;font-size:12px;display:flex;align-items:center;gap:9px;
  animation:si .25s ease;box-shadow:0 6px 28px rgba(0,0,0,.5);}
@keyframes si{from{opacity:0;transform:translateX(14px)}to{opacity:1;transform:none}}
.notif.ok{border-left:3px solid var(--acc);}
.notif.er{border-left:3px solid var(--acc3);}
.notif.in{border-left:3px solid var(--acc2);}

/* MISC */
.empty{padding:50px 20px;text-align:center;color:var(--txt3);}
.empty-ic{font-size:32px;margin-bottom:12px;}
.empty-t{font-size:13px;color:var(--txt2);}
.empty-s{font-size:11px;font-family:var(--mono);margin-top:4px;}
.spin{width:14px;height:14px;border:2px solid transparent;border-top-color:currentColor;
  border-radius:50%;animation:sp .6s linear infinite;display:inline-block;}
@keyframes sp{to{transform:rotate(360deg)}}
.re{display:flex;justify-content:flex-end;margin-bottom:14px;}
.hint{font-size:11px;color:var(--txt2);font-family:var(--mono);padding:9px 12px;
  background:var(--sur2);border-radius:5px;border-left:2px solid var(--acc2);margin-top:14px;}
::-webkit-scrollbar{width:4px;height:4px;}
::-webkit-scrollbar-thumb{background:var(--bdr2);border-radius:2px;}
</style>
</head>
<body>

<!-- LOGIN -->
<div id="login">
  <div class="lw">
    <div class="lb">
      <div class="li">FRP</div>
      <div class="ln">frpc<em>-panel</em></div>
    </div>
    <div class="ls">// BUILT-IN MANAGEMENT CONSOLE</div>
    <div class="fl"><label>用户名</label><input class="fi" id="lu" value="admin" autocomplete="username"></div>
    <div class="fl"><label>密码</label><input class="fi" id="lp" type="password" placeholder="••••••" autocomplete="current-password"></div>
    <button class="btn bp" onclick="login()">登录</button>
    <div class="emsg" id="lerr">用户名或密码错误</div>
    <div class="hint">默认凭据: admin / admin<br>凭据存储在 frpc-panel-auth.json（config 文件同目录）</div>
  </div>
</div>

<!-- APP -->
<div id="app">
  <aside class="side">
    <div class="sh">
      <div class="si">FRP</div>
      <div class="st">frpc<em>-panel</em></div>
    </div>
    <nav class="snav">
      <button class="ni on" onclick="go('ov')" id="n-ov"><span class="nic">⬡</span>概览</button>
      <button class="ni" onclick="go('px')" id="n-px"><span class="nic">⇄</span>代理管理</button>
      <button class="ni" onclick="go('cf')" id="n-cf"><span class="nic">◈</span>配置文件</button>
      <button class="ni" onclick="go('st')" id="n-st"><span class="nic">⚙</span>设置</button>
    </nav>
    <div class="sf">
      <div class="sp"><div class="dot" id="dot"></div><span id="sdot" style="font-size:11px;">连接中...</span></div>
      <button class="ni" onclick="logout()" style="color:var(--acc3);"><span class="nic">⎋</span>退出</button>
    </div>
  </aside>

  <main class="main">
    <div class="top">
      <div><div class="pt" id="ptl">概览</div><div class="ps" id="psl">// STATUS</div></div>
      <div style="display:flex;gap:7px;align-items:center;">
        <button class="bic" onclick="refresh()" title="刷新">↻</button>
        <span id="uname" style="font-size:11px;color:var(--txt2);font-family:var(--mono);
          padding:3px 9px;background:var(--sur2);border-radius:4px;border:1px solid var(--bdr);">admin</span>
      </div>
    </div>

    <div class="con">

      <!-- OVERVIEW -->
      <div class="sec on" id="sec-ov">
        <div class="cb2">
          <button class="cc cs" onclick="doReload()">↺ 重载配置</button>
          <button class="cc" onclick="fetchStatus()">● 刷新状态</button>
        </div>
        <div class="sg">
          <div class="sc c0"><div class="sl">代理总数</div><div class="sv" id="s0">—</div></div>
          <div class="sc c1"><div class="sl">运行中</div><div class="sv" id="s1">—</div></div>
          <div class="sc c2"><div class="sl">错误</div><div class="sv" id="s2">—</div></div>
          <div class="sc c3"><div class="sl">运行时间</div><div class="sv" id="s3" style="font-size:16px;">—</div></div>
        </div>
        <div class="card">
          <div class="ch"><div class="ct">⇄ 代理状态</div>
            <button class="bic" onclick="fetchStatus()" title="刷新">↻</button></div>
          <div id="ov-px"></div>
        </div>
      </div>

      <!-- PROXIES -->
      <div class="sec" id="sec-px">
        <div class="re"><button class="btn bp" style="width:auto;" onclick="openAdd()">+ 添加代理</button></div>
        <div class="card">
          <div class="ch"><div class="ct">⇄ 代理列表</div>
            <button class="bic" onclick="loadProxies()" title="刷新">↻</button></div>
          <div id="px-list"></div>
        </div>
      </div>

      <!-- CONFIG -->
      <div class="sec" id="sec-cf">
        <div class="card">
          <div class="ch">
            <div class="ct">◈ frpc.toml</div>
            <div style="display:flex;gap:7px;">
              <button class="btn bg bsm" onclick="loadCfg()">加载</button>
              <button class="btn bp bsm" onclick="saveCfg()">保存 &amp; 重载</button>
            </div>
          </div>
          <div class="cb">
            <textarea class="fta" id="cfgtxt" style="min-height:340px;" placeholder="点击「加载」读取当前配置内容..."></textarea>
            <div class="hint">修改配置后点击「保存 &amp; 重载」，frpc 会立即应用新配置（热重载，不断连接）。</div>
          </div>
        </div>
      </div>

      <!-- SETTINGS -->
      <div class="sec" id="sec-st">
        <div class="card">
          <div class="ch"><div class="ct">⚿ 修改面板登录密码</div></div>
          <div class="cb">
            <div class="g1"><div class="fl" style="margin:0;margin-bottom:12px;"><label>当前密码</label>
              <input class="fi" id="pc" type="password"></div></div>
            <div class="g2">
              <div class="fl" style="margin:0;"><label>新密码</label><input class="fi" id="pn" type="password"></div>
              <div class="fl" style="margin:0;"><label>确认新密码</label><input class="fi" id="pk" type="password"></div>
            </div>
            <div style="margin-top:12px;">
              <button class="btn bp" style="width:auto;" onclick="chpwd()">修改密码</button>
            </div>
            <div class="hint">面板密码独立于 frpc webServer 认证，存储在 frpc-panel-auth.json。</div>
          </div>
        </div>
      </div>

    </div>
  </main>
</div>

<!-- ADD/EDIT PROXY MODAL -->
<div class="ov" id="m-px">
  <div class="modal">
    <div class="mh"><div class="mtl" id="m-px-title">添加代理</div><button class="mcl" onclick="cm('m-px')">✕</button></div>
    <div class="mb">
      <div class="g1"><div class="fl" style="margin:0;margin-bottom:12px;"><label>代理名称</label>
        <input class="fi" id="pxn" placeholder="my-ssh"></div></div>
      <div class="g2">
        <div class="fl" style="margin:0;"><label>类型</label>
          <select class="fs" id="pxt" onchange="onType()">
            <option value="tcp">TCP</option><option value="udp">UDP</option>
            <option value="http">HTTP</option><option value="https">HTTPS</option>
          </select></div>
        <div class="fl" style="margin:0;"><label>本地 IP</label><input class="fi" id="pxli" value="127.0.0.1"></div>
      </div>
      <div class="g2">
        <div class="fl" style="margin:0;"><label>本地端口 (localPort)</label><input class="fi" id="pxlp" type="number" placeholder="22"></div>
        <div class="fl" style="margin:0;" id="pxrpw"><label>远程端口 (remotePort)</label><input class="fi" id="pxrp" type="number" placeholder="6000"></div>
      </div>
      <div id="pxhf" style="display:none;">
        <div class="g2">
          <div class="fl" style="margin:0;"><label>自定义域名 (customDomains)</label><input class="fi" id="pxcd" placeholder="example.com"></div>
          <div class="fl" style="margin:0;"><label>子域名 (subdomain)</label><input class="fi" id="pxsd" placeholder="myapp"></div>
        </div>
      </div>
      <div class="g2" style="margin-top:12px;">
        <div class="fl" style="margin:0;"><label>带宽限制 (bandwidthLimit)</label><input class="fi" id="pxbw" placeholder="如: 1MB"></div>
        <div style="display:flex;align-items:center;gap:9px;padding-top:24px;">
          <label class="tog"><input type="checkbox" id="pxenc"><span class="tsl"></span></label>
          <span style="font-size:12px;color:var(--txt2);">加密 (useEncryption)</span>
        </div>
      </div>
    </div>
    <div class="mft">
      <button class="btn bg" onclick="cm('m-px')">取消</button>
      <button class="btn bp" style="width:auto;" onclick="saveProxy()">保存</button>
    </div>
  </div>
</div>

<!-- CONFIRM MODAL -->
<div class="ov" id="m-cfm">
  <div class="modal" style="width:340px;">
    <div class="mh"><div class="mtl" id="cfmt">确认</div><button class="mcl" onclick="cm('m-cfm')">✕</button></div>
    <div class="mb"><p style="color:var(--txt2);font-size:13px;" id="cfmm"></p></div>
    <div class="mft">
      <button class="btn bg" onclick="cm('m-cfm')">取消</button>
      <button class="btn bd" style="width:auto;" onclick="cfmok()">确认删除</button>
    </div>
  </div>
</div>

<script>
// ── BASE PATH ─────────────────────────────────────────────────────────────────
// All API calls are relative to /panel/api/
const BASE = '/panel/api';

// ── STATE ─────────────────────────────────────────────────────────────────────
const S = { sec: 'ov', editPx: null, cfmCb: null };

// ── API ───────────────────────────────────────────────────────────────────────
const api = {
  async req(m, p, b) {
    const o = { method: m, credentials: 'include',
      headers: { 'Content-Type': p === '/config' && m === 'PUT' ? 'text/plain' : 'application/json' } };
    if (b !== undefined) o.body = typeof b === 'string' ? b : JSON.stringify(b);
    const r = await fetch(BASE + p, o);
    const ct = r.headers.get('content-type') || '';
    const d = ct.includes('json') ? await r.json() : await r.text();
    if (!r.ok) throw new Error((d && d.error) || r.statusText);
    return d;
  },
  get: p => api.req('GET', p),
  post: (p, b) => api.req('POST', p, b),
  put: (p, b) => api.req('PUT', p, b),
  del: p => api.req('DELETE', p),
};

// ── AUTH ──────────────────────────────────────────────────────────────────────
async function login() {
  const u = g('lu').value.trim(), p = g('lp').value;
  try {
    const r = await api.post('/login', { username: u, password: p });
    g('lerr').style.display = 'none';
    g('login').style.display = 'none';
    g('app').style.display = 'block';
    g('uname').textContent = r.username || u;
    initApp();
  } catch(e) { g('lerr').style.display = 'block'; }
}
g('lp').addEventListener('keydown', e => { if (e.key === 'Enter') login(); });

async function logout() {
  try { await api.post('/logout'); } catch(e) {}
  g('login').style.display = 'flex';
  g('app').style.display = 'none';
  g('lp').value = '';
}

async function chpwd() {
  const c = g('pc').value, n = g('pn').value, k = g('pk').value;
  if (n !== k) { notify('两次密码不一致', 'er'); return; }
  try {
    await api.post('/password', { current: c, new: n });
    g('pc').value = g('pn').value = g('pk').value = '';
    notify('密码已修改', 'ok');
  } catch(e) { notify(e.message, 'er'); }
}

// ── INIT ──────────────────────────────────────────────────────────────────────
function initApp() {
  fetchStatus();
  setInterval(fetchStatus, 15000);
}

// ── NAV ───────────────────────────────────────────────────────────────────────
const pages = {
  ov: ['概览', '// STATUS DASHBOARD'],
  px: ['代理管理', '// TUNNEL CONFIGURATION'],
  cf: ['配置文件', '// FRPC.TOML EDITOR'],
  st: ['设置', '// PANEL PREFERENCES'],
};
function go(id) {
  document.querySelectorAll('.sec').forEach(s => s.classList.remove('on'));
  document.querySelectorAll('.ni').forEach(n => n.classList.remove('on'));
  g('sec-' + id).classList.add('on');
  g('n-' + id).classList.add('on');
  g('ptl').textContent = pages[id][0];
  g('psl').textContent = pages[id][1];
  S.sec = id;
  if (id === 'px') loadProxies();
  if (id === 'cf') loadCfg();
  if (id === 'ov') fetchStatus();
}
function refresh() {
  if (S.sec === 'ov') fetchStatus();
  if (S.sec === 'px') loadProxies();
  if (S.sec === 'cf') loadCfg();
}

// ── STATUS ────────────────────────────────────────────────────────────────────
async function fetchStatus() {
  // Get frpc proxy status from our inline panel API
  try {
    const [st, info] = await Promise.all([api.get('/status'), api.get('/info')]);
    const all = [];
    for (const [, list] of Object.entries(st)) for (const p of list) all.push(p);
    g('s0').textContent = all.length;
    g('s1').textContent = all.filter(p => p.status === 'running').length;
    g('s2').textContent = all.filter(p => p.status !== 'running').length;
    g('s3').textContent = info.uptime || '—';

    const dot = g('dot'), sd = g('sdot');
    dot.className = 'dot r'; sd.textContent = '已连接';

    if (all.length === 0) {
      g('ov-px').innerHTML = '<div class="empty"><div class="empty-ic">⇄</div><div class="empty-t">暂无代理</div><div class="empty-s">在「代理管理」中添加隧道</div></div>';
      return;
    }
    const rows = all.map(p => {
      const tb = typeBadge(p.type);
      const sb = p.status === 'running'
        ? '<span class="bdg br">● 运行</span>'
        : '<span class="bdg be">✕ ' + esc(p.err || '错误') + '</span>';
      return '<tr><td><span class="pn">' + esc(p.name) + '</span></td><td>' + tb + '</td><td>' +
        esc(p.local_addr||'—') + '</td><td>' + esc(p.remote_addr||'—') + '</td><td>' + sb + '</td></tr>';
    }).join('');
    g('ov-px').innerHTML = '<table class="tbl"><thead><tr><th>名称</th><th>类型</th><th>本地</th><th>远程</th><th>状态</th></tr></thead><tbody>' + rows + '</tbody></table>';
  } catch(e) {
    g('dot').className = 'dot e'; g('sdot').textContent = '断开';
    g('s0').textContent = g('s1').textContent = g('s2').textContent = '—';
    g('ov-px').innerHTML = '<div class="empty"><div class="empty-ic">⚠</div><div class="empty-t">无法获取状态</div><div class="empty-s">' + esc(e.message) + '</div></div>';
  }
}

// ── RELOAD ───────────────────────────────────────────────────────────────────
async function doReload() {
  try { await api.post('/reload'); notify('配置已重载', 'ok'); setTimeout(fetchStatus, 800); }
  catch(e) { notify(e.message, 'er'); }
}

// ── PROXIES ───────────────────────────────────────────────────────────────────
async function loadProxies() {
  const el = g('px-list');
  el.innerHTML = '<div class="empty"><div class="spin"></div></div>';
  try {
    const r = await api.get('/proxies');
    const proxies = r.proxies || [];
    if (!proxies.length) {
      el.innerHTML = '<div class="empty"><div class="empty-ic">⇄</div><div class="empty-t">暂无代理</div><div class="empty-s">点击右上角按钮添加</div></div>';
      return;
    }
    const rows = proxies.map(p => {
      const cfg = p[p.type] || {};
      const remote = cfg.remotePort || (cfg.customDomains && cfg.customDomains[0]) || '—';
      return '<tr><td><span class="pn">' + esc(p.name) + '</span></td>' +
        '<td>' + typeBadge(p.type) + '</td>' +
        '<td>' + esc((cfg.localIP||'127.0.0.1') + ':' + (cfg.localPort||'—')) + '</td>' +
        '<td>' + esc(String(remote)) + '</td>' +
        '<td><div class="ag">' +
          '<button class="bic" onclick=\'openEdit(' + JSON.stringify(p) + ')\' title="编辑">✎</button>' +
          '<button class="bic" style="color:var(--acc3);" onclick="delProxy(\'' + esc(p.name) + '\')" title="删除">✕</button>' +
        '</div></td></tr>';
    }).join('');
    el.innerHTML = '<table class="tbl"><thead><tr><th>名称</th><th>类型</th><th>本地</th><th>远程</th><th>操作</th></tr></thead><tbody>' + rows + '</tbody></table>';
  } catch(e) {
    el.innerHTML = '<div class="empty"><div class="empty-ic">⚠</div><div class="empty-t">加载失败</div><div class="empty-s">' + esc(e.message) + '</div></div>';
  }
}

function openAdd() {
  S.editPx = null;
  g('m-px-title').textContent = '添加代理';
  g('pxn').value = ''; g('pxn').disabled = false;
  g('pxt').value = 'tcp'; g('pxli').value = '127.0.0.1';
  g('pxlp').value = ''; g('pxrp').value = '';
  g('pxcd').value = ''; g('pxsd').value = '';
  g('pxbw').value = ''; g('pxenc').checked = false;
  onType(); om('m-px');
}

function openEdit(p) {
  S.editPx = p.name;
  g('m-px-title').textContent = '编辑代理';
  const cfg = p[p.type] || {};
  g('pxn').value = p.name; g('pxn').disabled = true;
  g('pxt').value = p.type || 'tcp'; g('pxli').value = cfg.localIP || '127.0.0.1';
  g('pxlp').value = cfg.localPort || ''; g('pxrp').value = cfg.remotePort || '';
  g('pxcd').value = (cfg.customDomains && cfg.customDomains[0]) || '';
  g('pxsd').value = cfg.subdomain || '';
  g('pxbw').value = cfg.bandwidthLimit || '';
  g('pxenc').checked = !!(cfg.transport && cfg.transport.useEncryption);
  onType(); om('m-px');
}

function onType() {
  const t = g('pxt').value;
  const h = t === 'http' || t === 'https';
  g('pxhf').style.display = h ? 'block' : 'none';
  g('pxrpw').style.display = h ? 'none' : 'block';
}

async function saveProxy() {
  const name = g('pxn').value.trim();
  const type = g('pxt').value;
  const localIP = g('pxli').value.trim() || '127.0.0.1';
  const localPort = parseInt(g('pxlp').value) || 0;
  const isHttp = type === 'http' || type === 'https';
  const remotePort = parseInt(g('pxrp').value) || 0;
  if (!name) { notify('请填写代理名称', 'er'); return; }
  if (!localPort) { notify('请填写本地端口', 'er'); return; }
  if (!isHttp && !remotePort) { notify('请填写远程端口', 'er'); return; }

  // Build the ProxyDefinition the same way the store API expects
  const typeCfg = { localIP, localPort };
  if (!isHttp) typeCfg.remotePort = remotePort;
  const cd = g('pxcd').value.trim();
  if (isHttp && cd) typeCfg.customDomains = [cd];
  const sd = g('pxsd').value.trim();
  if (isHttp && sd) typeCfg.subdomain = sd;
  const bw = g('pxbw').value.trim();
  if (bw) typeCfg.bandwidthLimit = bw;
  if (g('pxenc').checked) typeCfg.transport = { useEncryption: true };

  const payload = { name, type, [type]: typeCfg };

  try {
    if (S.editPx) {
      await api.put('/proxies/' + encodeURIComponent(S.editPx), payload);
      notify('代理已更新，正在重载...', 'ok');
    } else {
      await api.post('/proxies', payload);
      notify('代理已添加，正在重载...', 'ok');
    }
    cm('m-px'); g('pxn').disabled = false;
    await api.post('/reload');
    await loadProxies();
    await fetchStatus();
  } catch(e) { notify(e.message, 'er'); }
}

function delProxy(name) {
  cfm('删除代理', '确定要删除代理 "' + name + '" 吗？删除后自动重载配置。', async () => {
    try {
      await api.del('/proxies/' + encodeURIComponent(name));
      await api.post('/reload');
      notify('已删除', 'ok');
      await loadProxies();
    } catch(e) { notify(e.message, 'er'); }
  });
}

// ── CONFIG ────────────────────────────────────────────────────────────────────
async function loadCfg() {
  try { g('cfgtxt').value = await api.get('/config'); }
  catch(e) { notify('加载失败: ' + e.message, 'er'); }
}

async function saveCfg() {
  try {
    await api.put('/config', g('cfgtxt').value);
    await api.post('/reload');
    notify('配置已保存并重载', 'ok');
    await fetchStatus();
  } catch(e) { notify(e.message, 'er'); }
}

// ── MODAL ─────────────────────────────────────────────────────────────────────
function om(id) { g(id).classList.add('on'); }
function cm(id) {
  g(id).classList.remove('on');
  if (id === 'm-px') { S.editPx = null; g('pxn').disabled = false; }
}
document.querySelectorAll('.ov').forEach(el => el.addEventListener('click', e => { if (e.target === el) cm(el.id); }));

function cfm(title, msg, cb) {
  g('cfmt').textContent = title; g('cfmm').textContent = msg;
  S.cfmCb = cb; om('m-cfm');
}
function cfmok() { cm('m-cfm'); if (S.cfmCb) S.cfmCb(); }

// ── NOTIFY ────────────────────────────────────────────────────────────────────
function notify(msg, t = 'in') {
  const el = document.createElement('div');
  el.className = 'notif ' + t;
  el.innerHTML = ({ok:'✓',er:'✕',in:'ℹ'}[t]||'ℹ') + ' ' + esc(msg);
  document.body.appendChild(el);
  setTimeout(() => el.remove(), 3200);
}

// ── UTILS ─────────────────────────────────────────────────────────────────────
function g(id) { return document.getElementById(id); }
function esc(s) { return String(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;'); }
function typeBadge(t) {
  const m = {tcp:'bt',udp:'bu',http:'bh',https:'bh'};
  return '<span class="bdg ' + (m[t]||'bt') + '">' + (t||'tcp').toUpperCase() + '</span>';
}

// ── BOOT ──────────────────────────────────────────────────────────────────────
(async () => {
  try {
    const r = await api.get('/ping');
    if (r.authenticated) {
      g('login').style.display = 'none';
      g('app').style.display = 'block';
      initApp();
    }
  } catch(e) { /* show login */ }
})();
</script>
</body>
</html>`
