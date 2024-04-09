import { loginByUsernamePassword } from '../login_ui.js';
import { createOrg } from '../setup.js';
import { createMachine, addMachinePat } from '../user.js';
import { removeOrg } from '../teardown.js';
import { userinfo } from '../oidc.js';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../config.js';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin);
  console.info("setup: admin signed in");

  const org = await createOrg(tokens.accessToken);
  console.info(`setup: org (${org.organizationId}) created`);

  let machines = Array.from({length: MaxVUs()}, (_, i) => {
    return createMachine(`zitachine-${i}`, org, tokens.accessToken);
  });
  machines = await Promise.all(machines);
  machines = machines.map((machine) => {
    return {userId: machine.userId, loginName: machine.loginNames[0]};
  });
  console.info(`setup: ${machines.length} machines created`);

  let pats = machines.map((machine) => {
    return addMachinePat(machine.userId, org, tokens.accessToken);
  });
  pats = await Promise.all(pats);
  machines = pats.map((pat, i) => {
    return {userId: machines[i].userId, loginName: machines[i].loginName, pat: pat.token};
  });
  console.info(`setup: Pats added`);

  return {tokens, machines: machines, org};
}

const humanPasswordLoginTrend = new Trend('machine_pat_login_duration', true);
export default function (data) {
  const token = userinfo(data.machines[__VU-1].pat);
}

export function teardown(data) {
  removeOrg(data.org, data.tokens.accessToken);
  console.info('teardown: org removed')
}

