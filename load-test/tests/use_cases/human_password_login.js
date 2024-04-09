import { loginByUsernamePassword } from '../login_ui.js';
import { createOrg} from '../setup.js';
import { createHuman } from '../user.js';
import { removeOrg } from '../teardown.js';
import { userinfo } from '../oidc.js';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../config.js';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin);
  console.log("setup: admin signed in");
  
  const org = await createOrg(tokens.accessToken);
  console.log(`setup: org (${org.organizationId}) created`);

  let humans = Array.from({length: MaxVUs()}, (_, i) => {
    return createHuman(`zitizen-${i}`, org, tokens.accessToken);
  });
  humans = await Promise.all(humans);
  humans = humans.map((user, i) => {
    return {userId: user.userId, loginName: user.loginNames[0], password: 'Password1!'};
  })
  console.log(`setup: ${humans.length} users created`);
  return {tokens, users: humans, org};
}

const humanPasswordLoginTrend = new Trend('human_password_login_duration', true);
export default function (data) {
  const start = new Date();
  const token = loginByUsernamePassword(data.users[__VU-1]);
  userinfo(token.accessToken);

  humanPasswordLoginTrend.add(new Date() - start);
}

export function teardown(data) {
  removeOrg(data.org, data.tokens.accessToken);
}

