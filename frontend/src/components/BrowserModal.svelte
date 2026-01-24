<script>
  import { createEventDispatcher } from 'svelte';
  import { SelectFile, SelectDirectory } from '../../wailsjs/go/main/App';
  
  const dispatch = createEventDispatcher();

  export let show = false;
  export let instance = null;

  let name = '';
  let path = '';
  let userDataDir = '';
  let argsString = '';

  let lastShow = false;
  let lastInstanceId = null;

  $: if (show && (show !== lastShow || (instance && instance.id !== lastInstanceId))) {
    // 仅在模态框由关闭转开启，或者切换了不同的实例时，才同步数据
    if (instance) {
      name = instance.name || '';
      path = instance.path || '';
      userDataDir = instance.userDataDir || '';
      argsString = formatArgs(instance.args || []);
      lastInstanceId = instance.id;
    } else {
      name = '';
      path = '';
      userDataDir = '';
      argsString = '';
      lastInstanceId = null;
    }
    lastShow = show;
  } else if (!show && lastShow) {
    // 关闭时更新状态
    lastShow = false;
  }

  async function handleBrowseFile() {
    const result = await SelectFile();
    if (result) path = result;
  }

  async function handleBrowseDir() {
    const result = await SelectDirectory();
    if (result) userDataDir = result;
  }

  function formatArgs(args) {
    if (!args || args.length === 0) return '';
    return args.map(arg => {
        if (arg.includes(' ') && !arg.startsWith('"') && !arg.startsWith("'")) {
            // Check if it's a key=value pair
            if (arg.includes('=')) {
                const parts = arg.split('=');
                if (parts.length === 2) {
                     return `${parts[0]}="${parts[1]}"`;
                }
            }
            return `"${arg}"`;
        }
        return arg;
    }).join(' ');
  }

  function parseArgs(str) {
    const args = [];
    let current = '';
    let quote = null;
    let escaped = false;

    for (let i = 0; i < str.length; i++) {
        const char = str[i];
        
        if (escaped) {
            current += char;
            escaped = false;
            continue;
        }

        if (char === '\\') {
            escaped = true;
            continue;
        }

        if (quote) {
            if (char === quote) {
                quote = null;
            } else {
                current += char;
            }
        } else {
            if (char === '"' || char === "'") {
                quote = char;
            } else if (char === ' ') {
                if (current.length > 0) {
                    args.push(current);
                    current = '';
                }
            } else {
                current += char;
            }
        }
    }
    
    if (current.length > 0) {
        args.push(current);
    }
    
    return args;
  }

  function handleSave() {
    if (!name || !path) {
      alert("请填写名称和浏览器路径");
      return;
    }
    const args = parseArgs(argsString);
    dispatch('save', {
      id: instance ? instance.id : null,
      name,
      path,
      userDataDir,
      args
    });
  }

  function handleClose() {
    dispatch('close');
  }
</script>

{#if show}
  <div class="modal-overlay" on:click|self={handleClose}>
    <div class="modal">
      <h2 style="margin-bottom: 1.5rem; font-size: 1.25rem;">{instance ? '编辑实例' : '添加实例'}</h2>
      
      <div class="form-group">
        <label>名称</label>
        <input class="form-control" type="text" bind:value={name} placeholder="例如：账号 A" />
      </div>

      <div class="form-group">
        <label>浏览器路径</label>
        <div style="display: flex; gap: 0.5rem;">
          <input class="form-control" type="text" bind:value={path} placeholder="可执行文件路径" />
          <button class="btn" style="padding: 0 1rem; flex-shrink: 0;" on:click={handleBrowseFile}>浏览...</button>
        </div>
      </div>

      <div class="form-group">
        <label>用户数据目录 (隔离指纹)</label>
        <div style="display: flex; gap: 0.5rem;">
          <input class="form-control" type="text" bind:value={userDataDir} placeholder="配置保存目录" />
          <button class="btn" style="padding: 0 1rem; flex-shrink: 0;" on:click={handleBrowseDir}>选择...</button>
        </div>
      </div>

      <div class="form-group">
        <label>启动参数 (空格分隔)</label>
        <textarea class="form-control" rows="3" bind:value={argsString} placeholder="--incognito --proxy-server=..." style="resize: vertical;"></textarea>
      </div>

      <div style="display: flex; gap: 1rem; margin-top: 2rem; justify-content: flex-end;">
        <button class="btn" on:click={handleClose}>取消</button>
        <button class="btn btn-primary" on:click={handleSave}>保存配置</button>
      </div>
    </div>
  </div>
{/if}

<style>
  textarea.form-control {
    font-family: monospace;
    font-size: 0.75rem;
  }
</style>
