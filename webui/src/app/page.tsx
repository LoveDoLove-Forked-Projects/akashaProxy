'use client'

import Container from '@mui/material/Container';
import Chip from '@mui/material/Chip';
import InfoLine from '@/component/InfoLine';
import Stack from '@mui/material/Stack';
import Fab from '@mui/material/Fab';
import Button from '@mui/material/Button';
import { blueGrey, grey, lightBlue } from '@mui/material/colors';
import RotateRight from '@mui/icons-material/RotateRight';
import DoneAll from '@mui/icons-material/DoneAll';
import ClearIcon from '@mui/icons-material/Clear';
import Settings from '@mui/icons-material/Settings';
import GpsFixedIcon from '@mui/icons-material/GpsFixed';
import SpeedIcon from '@mui/icons-material/Speed';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import useMediaQuery from '@mui/material/useMediaQuery';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import LoadingButton from '@mui/lab/LoadingButton';
import { exec, toast } from 'kernelsu';
import { useEffect, useState, useMemo } from 'react';
import { MIHOMO_PATH } from './consts';
import yaml from 'js-yaml';
import Card from '@mui/material/Card';
import CssBaseline from '@mui/material/CssBaseline';


interface MihomoInfo {
  version: string | null,
  versions: string | null,
  daemon: number | null,
  webui: string | null,
  log: string | null,
  iptest: string | null,
  ipspeed: string | null,
  loading: boolean,
}
function MihomoCard({ info, setMihomoInfo, dark }: { info: MihomoInfo, setMihomoInfo: (callback: (info: MihomoInfo) => MihomoInfo) => void, dark: boolean }) {
  let bigButtonStyle: React.CSSProperties = {
    width: "100%",
    padding: "20px",
    textTransform: "none",
    fontSize: "18px",
    marginTop: "1rem",
    justifyContent: "start",
    borderRadius: "10px",
  };

  const bigButtonTheme = useMemo(() => createTheme({
    palette: dark ?
      {
        success: {
          main: lightBlue[900],
        },
        warning: {
          main: blueGrey[500],
        },
        info: {
          main: grey[800],
          dark: grey[900],
        },
        mode: 'dark',
      }
      : {
        success: {
          main: lightBlue[900],
        },
        warning: {
          main: blueGrey[500],
        },
        info: {
          main: grey[50],
          dark: grey[100],
        },
      },
  }), [dark]);

  return (
    <>
      <ThemeProvider theme={bigButtonTheme}>
        <LoadingButton variant="contained" style={{
          width: "100%",
          padding: "24px 20px",
          textTransform: "none",
          fontSize: "20px",
          marginTop: "1rem",
          justifyContent: "start",
          borderRadius: "10px",
        }}
          loading={info.loading}
          startIcon={info.daemon != null ? (<DoneAll />) : (<ClearIcon />)}
          color={info.daemon != null ? "success" : "warning"}
          onClick={() => info.daemon != null ? stopMihomo(setMihomoInfo) : startMihomo(setMihomoInfo)}
        >{info.daemon != null ? "Mihomo 运行正常" : "Mihomo 已停止"}</LoadingButton>

        {info.daemon != null && (
          <Button variant="contained"
            style={bigButtonStyle}
            startIcon={<Settings />}
            color="info"
            onClick={() => window.location.href = info.webui ?? ""}
          >网页面板</Button>
        )}

        <Button variant="contained"
          style={bigButtonStyle}
          startIcon={<GpsFixedIcon />}
          color="info"
          onClick={() => window.location.href = info.iptest ?? "https://ip.skk.moe"}
        >IP检查</Button>

        <Button variant="contained"
          style={bigButtonStyle}
          startIcon={<SpeedIcon />}
          color="info"
          onClick={() => window.location.href = info.ipspeed ?? "https://fast.com"}
        >网络测速</Button>
      </ThemeProvider>

      <Card style={{ padding: '16px', backgroundColor: dark ? '#303030' : '#fafafc', marginTop: "2rem" }}>
        <Stack spacing={1}>
          <InfoLine title="Mihomo版本">
            {(info.version != null && <Chip label={info.version} color="primary" variant="outlined" size='small' />)}
            {(info.version == null && <Chip label="未知" color="warning" variant="outlined" size='small' />)}

          </InfoLine>
          <InfoLine title="模块版本">
            {(info.versions != null && <Chip label={info.versions} color="primary" variant="outlined" size='small' />)}
            {(info.versions == null && <Chip label="未知" color="warning" variant="outlined" size='small' />)}
          </InfoLine>

          {(info.log != null &&
            <>
              <Divider />
              <div>
                <p style={{ marginBottom: '1ex' }}>日志</p>
                <pre style={{
                  overflowX: "scroll",
                  borderRadius: '5px',
                  padding: '8px',
                  backgroundColor: dark ? "#141414" : "#eaeaec",
                  lineHeight: '1.3',
                  fontSize: '12px',
                  fontFamily: "monospace",
                  display: "block",
                  whiteSpace: 'pre',
                  userSelect: 'text',
                }}>{info.log}</pre>
              </div>
            </>
          )}
        </Stack>
      </Card>
    </>
  )
}

async function startMihomo(setMihomoInfo: (callback: (info: MihomoInfo) => MihomoInfo) => void) {
  setMihomoInfo(info => ({ ...info, loading: true }));
  try {
    // we need change cgroup before start mihomo, otherwise mihomo will be killed when kernelsu manager is be killed
    let process = await exec(`su -c "${MIHOMO_PATH}/tools/start.sh"`);
    if (process.errno != 0) {
      throw 'Failed: Exit code ' + process.errno + '\noutput: ' + process.stdout + '\nstderr: ' + process.stderr;
    }
    await updateInfo(setMihomoInfo);
  } catch (err) {
    setMihomoInfo(info => ({ ...info, loading: false }));
    console.error(err);
    toast("" + err);
  }
}

async function stopMihomo(setMihomoInfo: (callback: (info: MihomoInfo) => MihomoInfo) => void) {
  setMihomoInfo(info => ({ ...info, loading: true }));
  try {
    let process = await exec(MIHOMO_PATH + '/tools/stop.sh');
    if (process.errno != 0) {
      throw 'Failed: Exit code ' + process.errno + '\noutput: ' + process.stdout + '\nstderr: ' + process.stderr;
    }
    await updateInfo(setMihomoInfo);
  } catch (err) {
    setMihomoInfo(info => ({ ...info, loading: false }));
    console.error(err);
    toast("" + err);
  }
}

async function deleteCache() {
  try {
    let process = await exec(MIHOMO_PATH + '/tools/DeleteCache.sh');
    if (process.errno != 0) {
      throw 'Failed: Exit code ' + process.errno + '\noutput: ' + process.stdout + '\nstderr: ' + process.stderr;
    }
    toast('Cache deleted');
  } catch (err) {
    console.error(err);
    toast("" + err);
  }
}

async function updateInfo(setMihomoInfo: (callback: (info: MihomoInfo) => MihomoInfo) => void) {
  let resultInfo: MihomoInfo = { version: null, versions: null, daemon: null, webui: null, iptest: null, ipspeed: null, log: null, loading: false };
  let running = false;
  // Get mihomo file name
  let fileName = null;
  try {
    let cmd = `source ${MIHOMO_PATH}/clash.config && printf "%s" $bin_name`;
    let process = await exec(cmd);
    if (process.errno != 0) {
      throw 'Failed to execute `' + cmd + '`: Exit code ' + process.errno;
    }
    fileName = process.stdout;
  } catch (err) {
    console.error(err);
  }
  // Get clash version
  if (fileName) {
    try {
      let cmd = `"${MIHOMO_PATH}/kernel/${fileName}" -v`;
      let process = await exec(cmd);
      if (process.errno != 0) {
        throw 'Failed to execute `' + cmd + '`: Exit code ' + process.errno;
      }
      let version = process.stdout;
      let versionMatch = version.match(/^(Mihomo|Mihomo Meta) ([^ ]+)/);
      if (versionMatch == null) {
        throw 'Failed to parse version from `' + version + '`';
      }
      version = versionMatch[2];
      resultInfo.version = version;
    } catch (err) {
      console.error(err);
    }
    // get daemon pid
    try {
      let cmd = `source ${MIHOMO_PATH}/clash.config && printf "%s" $(busybox pidof $kernel_bin)`;
      let process = await exec(cmd);
      if (process.errno != 0) {
        throw 'Failed to execute `' + cmd + '`: Exit code ' + process.errno;
      }
      let pid = process.stdout.trim();
      // check if pid is running
      let getExeProcess = await exec(`test $(realpath /proc/${pid}/exe) == $(realpath "${MIHOMO_PATH}/kernel/${fileName}") || exit 1`);
      if (getExeProcess.errno == 0) {
        running = true;
        resultInfo.daemon = parseInt(pid);
      }
      
    } catch (err) {
      console.error(err);
    }
  }

  if (running) {
    // Get webui address
    try {
      let configYamlProcess = await exec('cat ' + MIHOMO_PATH + '/run/config.yaml');
      if (configYamlProcess.errno != 0) {
        throw 'Failed to execute `cat ' + MIHOMO_PATH + '/run/webui.addr`: Exit code ' + configYamlProcess.errno;
      }
      let configYaml = configYamlProcess.stdout;
      let config = yaml.load(configYaml) as {
        'external-controller': string | undefined,
        'secret': string | undefined,
      };
      if (config['external-controller'] == undefined) throw 'external-controller not found in config.yaml';
      let port = config['external-controller'].split(':')[1];
      let url = 'http://127.0.0.1:' + port + '/ui/?hostname=127.0.0.1&port=' + port;
      if (config['secret'] != undefined) url += '&secret=' + config['secret'];
      resultInfo.webui = url;
    } catch (err) {
      console.error(err);
    }
  }
  // get iptest
  try {
    let cmd = `source ${MIHOMO_PATH}/clash.config && printf "%s" $iptest`;
    let process = await exec(cmd);
    if (process.errno != 0) {
      throw 'Failed to execute `' + cmd + '`: Exit code ' + process.errno;
    }
    resultInfo.iptest = process.stdout;
  } catch (err) {
    console.error(err);
  }
  // get ipspeed
  try {
    let cmd = `source ${MIHOMO_PATH}/clash.config && printf "%s" $ipspeed`;
    let process = await exec(cmd);
    if (process.errno != 0) {
      throw 'Failed to execute `' + cmd + '`: Exit code ' + process.errno;
    }
    resultInfo.ipspeed = process.stdout;
  } catch (err) {
    console.error(err);
  }
  // get module version
  try {
    let cmd = `source ${MIHOMO_PATH}/clash.config && printf "%s" $Module_version`;
    let process = await exec(cmd);
    if (process.errno != 0) {
      throw 'Failed to execute `' + cmd + '`: Exit code ' + process.errno;
    }
    resultInfo.versions = process.stdout;
  } catch (err) {
    console.error(err);
  }

  // get logs
  try {
    let logProcess = await exec('tail -n 100 ' + MIHOMO_PATH + '/run/run.logs');
    if (logProcess.errno != 0) {
      throw 'Failed to execute `tail -n 100 ' + MIHOMO_PATH + '/run/run.logs`: Exit code ' + logProcess.errno;
    }
    resultInfo.log = logProcess.stdout;
  } catch (err) {
    console.error(err);
  }
  setMihomoInfo(() => resultInfo);
}

function saveMihomoInfo(mihomoInfo: MihomoInfo) {
  if (!mihomoInfo.loading && typeof window !== "undefined") {
    window.sessionStorage.setItem('mihomoInfo', JSON.stringify(mihomoInfo));
  }
}

function loadMihomoInfo(): MihomoInfo {
  const defaultInfo: MihomoInfo = { version: null, versions: null, daemon: null, webui: null, log: null, iptest: null, ipspeed: null, loading: true };
  if (typeof window !== "undefined") {
    let mihomoInfo = window.sessionStorage.getItem('mihomoInfo');
    if (mihomoInfo == null) {
      return defaultInfo;
    }
    try {
      return JSON.parse(mihomoInfo);
    } catch (err) {
      console.error(err);
      return defaultInfo;
    }
  } else {
    return defaultInfo;
  }
}

export default function Home() {
  let [mihomoInfo, setMihomoInfo] = useState<MihomoInfo>({ version: null, versions: null, daemon: null, webui: null, log: null, iptest: null, ipspeed: null, loading: true });
  useEffect(() => saveMihomoInfo(mihomoInfo), [mihomoInfo]);
  useEffect(() => {
    loadMihomoInfo();
    updateInfo(setMihomoInfo);
  }, [])
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  return (
    <ThemeProvider theme={createTheme({
      palette: { mode: prefersDarkMode ? 'dark' : 'light' },
    })}>
      <CssBaseline >

        <div style={{ padding: '1em', userSelect: 'none' }}>
          <Typography variant="h5" style={{ marginBottom: '1em' }}>
            akashaProxy
          </Typography>
          <Container maxWidth="md" style={{ paddingLeft: '0px', paddingRight: '0px' }}>
            <MihomoCard info={mihomoInfo} setMihomoInfo={setMihomoInfo} dark={prefersDarkMode} />
          </Container>

          <Fab size="small" color="success" style={{ right: '1em', bottom: '1em', zIndex: 999, position: 'fixed' }} onClick={() => updateInfo(setMihomoInfo)}>
            <RotateRight />
          </Fab>
        </div>
      </CssBaseline>
    </ThemeProvider>
  )
}
