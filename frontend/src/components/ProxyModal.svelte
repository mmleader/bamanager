<script>
  import { createEventDispatcher } from 'svelte';
  
  const dispatch = createEventDispatcher();

  export let show = false;
  export let proxy = null;

  let name = '';
  let protocol = 'socks5';
  let host = '';
  let port = 1080;
  let username = '';
  let password = '';

  let lastShow = false;
  let lastProxyId = null;

  $: if (show && (show !== lastShow || (proxy && proxy.id !== lastProxyId))) {
    if (proxy) {
      name = proxy.name || '';
      protocol = proxy.protocol || 'socks5';
      host = proxy.host || '';
      port = proxy.port || 1080;
      username = proxy.username || '';
      password = proxy.password || '';
      lastProxyId = proxy.id;
    } else {
      name = '';
      protocol = 'socks5';
      host = '';
      port = 1080;
      username = '';
      password = '';
      lastProxyId = null;
    }
    lastShow = show;
  } else if (!show && lastShow) {
    lastShow = false;
  }

  function handleSave() {
    if (!name || !host || !port) {
      alert("请填写代理名称、地址和端口");
      return;
    }
    dispatch('save', {
      id: proxy ? proxy.id : null,
      name,
      protocol,
      host,
      port: parseInt(port),
      username,
      password
    });
  }

  function handleClose() {
    dispatch('close');
  }
</script>

{#if show}
  <div class="modal-overlay" on:click|self={handleClose}>
    <div class="modal">
      <h2 style="margin-bottom: 1.5rem; font-size: 1.25rem;">{proxy ? '编辑代理' : '添加代理'}</h2>
      
      <div class="form-group">
        <label>代理名称</label>
        <input class="form-control" type="text" bind:value={name} placeholder="例如：日本节点" />
      </div>

      <div class="form-group">
        <label>协议</label>
        <select class="form-control" bind:value={protocol}>
          <option value="socks5">SOCKS5</option>
          <option value="http">HTTP</option>
          <option value="https">HTTPS</option>
        </select>
      </div>

      <div class="form-group">
        <label>地址</label>
        <div style="display: flex; gap: 0.5rem;">
          <input class="form-control" type="text" bind:value={host} placeholder="代理服务器地址" style="flex: 1;" />
          <input class="form-control" type="number" bind:value={port} placeholder="端口" style="width: 100px;" />
        </div>
      </div>

      <div class="form-group">
        <label>认证信息（可选）</label>
        <div style="display: flex; gap: 0.5rem;">
          <input class="form-control" type="text" bind:value={username} placeholder="用户名" />
          <input class="form-control" type="password" bind:value={password} placeholder="密码" />
        </div>
      </div>

      <div style="display: flex; gap: 1rem; margin-top: 2rem; justify-content: flex-end;">
        <button class="btn" on:click={handleClose}>取消</button>
        <button class="btn btn-primary" on:click={handleSave}>保存代理</button>
      </div>
    </div>
  </div>
{/if}

<style>
  select.form-control {
    appearance: auto;
    cursor: pointer;
  }
</style>
