import React, { useEffect, useState } from "react";
import Link from "next/link";
import algosdk from "algosdk";
import hljs from "highlight.js";
import TabUnstyled from "@mui/base/TabUnstyled";
import TabsUnstyled from "@mui/base/TabsUnstyled";
import TabsListUnstyled from "@mui/base/TabsListUnstyled";
import TabPanelUnstyled from "@mui/base/TabPanelUnstyled";

import { TransactionResponse } from "../../types/apiResponseTypes";
import styles from "../../pages/tx/TransactionDetails.module.scss";
import blockStyles from "../../pages/block/Block.module.scss";
import { getAppTEAL, getClearStateTEAL } from "../../utils/api";
import { prettyPrintTEAL } from "../../utils/stringUtils";
import {
  algodAddr,
  algodProtocol,
  algodToken,
  isLocal,
} from "../../utils/constants";
import Copyable from "../copyable/Copyable";

const ApplicationTransactionInfo = ({ tx }: { tx: TransactionResponse }) => {
  const appTx = tx["application-transaction"];
  const [disassembledApp, setDisassembledApp] = useState<string>();
  const [disassembledClearStateProgram, setDisassembledClearStateProgram] =
    useState<string>();

  useEffect(() => {
    if (isLocal && algodToken && algodProtocol && algodAddr && tx && appTx) {
      if (appTx["approval-program"]) {
        getAppTEAL(tx)
          .then((result) => {
            if (
              result &&
              result.txns &&
              result.txns[0] &&
              result.txns[0].disassembly
            ) {
              const disassembledResult = prettyPrintTEAL(
                result.txns[0].disassembly
              );
              setDisassembledApp(disassembledResult);
            }
          })
          .catch((error) => {
            console.error("App disassembly error: ", error);
          });
      }
      if (appTx["clear-state-program"]) {
        getClearStateTEAL(tx)
          .then((result) => {
            if (
              result &&
              result.txns &&
              result.txns[0] &&
              result.txns[0].disassembly
            ) {
              const disassembledResult = prettyPrintTEAL(
                result.txns[0].disassembly
              );
              setDisassembledClearStateProgram(disassembledResult);
            }
          })
          .catch((error) => {
            console.error("App disassembly error: ", error);
          });
      }
    }
  }, [tx, appTx]);

  if (!appTx) return null;

  return (
    <div>
      <h4>Application Transaction Information</h4>
      <div className={blockStyles["block-table"]}>
        <table cellSpacing="0">
          <tbody>
            <tr>
              <td>Application ID</td>
              <td>{appTx["application-id"]}</td>
            </tr>
            <tr>
              <td className={styles["valign-top-identifier"]}>Accounts</td>
              <td className={styles["multiline-details"]}>
                {appTx.accounts && appTx.accounts.length > 0
                  ? appTx.accounts.map((ac) => (
                      <Link href={`/address/${ac}`} key={ac}>
                        {ac}
                      </Link>
                    ))
                  : "N/A"}
              </td>
            </tr>
            {appTx["application-args"] && appTx["application-args"].length > 0 && (
              <tr className={styles["valign-top-identifier"]}>
                <td>Arguments</td>
                <td className={styles["multiline-details"]}>
                  <div className={styles["inner-table-wrapper"]}>
                    <table className={styles["inner-table"]}>
                      <thead>
                        <tr>
                          <td>Base64</td>
                          <td>ASCII</td>
                          <td>UInt64</td>
                        </tr>
                      </thead>
                      <tbody>
                        {appTx["application-args"].map((appArg) => (
                          <tr key={appArg}>
                            <td>{appArg}</td>
                            <td>
                              {Buffer.from(appArg, "base64").toString("ascii")}
                            </td>
                            <td>
                              {Buffer.from(appArg, "base64").length <= 8
                                ? algosdk.decodeUint64(
                                    Buffer.from(appArg, "base64"),
                                    "mixed"
                                  )
                                : "N/A"}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </td>
              </tr>
            )}
            {appTx["foreign-apps"] && appTx["foreign-apps"].length > 0 && (
              <tr>
                <td>Foreign Apps</td>
                <td className={styles["multiline-details"]}>
                  {appTx["foreign-apps"].map((app) => (
                    <p key={app}>{app}</p>
                  ))}
                </td>
              </tr>
            )}
            {appTx["foreign-assets"] && appTx["foreign-assets"].length > 0 && (
              <tr>
                <td>Foreign Assets</td>
                <td className={styles["multiline-details"]}>
                  {appTx["foreign-assets"].map((asa) => (
                    <p key={asa}>{asa}</p>
                  ))}
                </td>
              </tr>
            )}
            <tr>
              <td>On Completion</td>
              <td>{appTx["on-completion"]}</td>
            </tr>
            {tx["created-application-index"] && (
              <tr>
                <td>Created Application Index</td>
                <td>{tx["created-application-index"]}</td>
              </tr>
            )}
            {appTx["global-state-schema"] && (
              <tr>
                <td>Global State Schema</td>
                <td className={styles["multiline-details"]}>
                  <p>
                    Number of byte-slice:{" "}
                    {appTx["global-state-schema"]["num-byte-slice"]}
                  </p>
                  <p>
                    Number of uint: {appTx["global-state-schema"]["num-uint"]}
                  </p>
                </td>
              </tr>
            )}
            {appTx["local-state-schema"] && (
              <tr>
                <td>Local State Schema</td>
                <td className={styles["multiline-details"]}>
                  <p>
                    Number of byte-slice:{" "}
                    {appTx["local-state-schema"]["num-byte-slice"]}
                  </p>
                  <p>
                    Number of uint: {appTx["local-state-schema"]["num-uint"]}
                  </p>
                </td>
              </tr>
            )}
            {appTx["approval-program"] && (
              <tr>
                <td className={styles["valign-top-identifier"]}>
                  Approval Program
                </td>
                <td>
                  {disassembledApp ? (
                    <TabsUnstyled defaultValue={0}>
                      <TabsListUnstyled className={styles.tabs}>
                        <TabUnstyled>TEAL</TabUnstyled>
                        <TabUnstyled>Base64</TabUnstyled>
                      </TabsListUnstyled>
                      <TabPanelUnstyled value={0}>
                        <Copyable copyableText={disassembledApp}>
                          <pre className={`${styles["teal-box"]} hljs`}>
                            <code
                              className="language-lua"
                              dangerouslySetInnerHTML={{
                                __html: hljs.highlight(disassembledApp, {
                                  language: "lua",
                                }).value,
                              }}
                            />
                          </pre>
                        </Copyable>
                      </TabPanelUnstyled>
                      <TabPanelUnstyled value={1}>
                        <Copyable copyableText={appTx["approval-program"]} />
                      </TabPanelUnstyled>
                    </TabsUnstyled>
                  ) : (
                    <Copyable copyableText={appTx["approval-program"]} />
                  )}
                </td>
              </tr>
            )}
            {appTx["clear-state-program"] && (
              <tr>
                <td className={styles["valign-top-identifier"]}>
                  Clear State Program
                </td>
                <td>
                  {disassembledClearStateProgram ? (
                    <TabsUnstyled defaultValue={0}>
                      <TabsListUnstyled className={styles.tabs}>
                        <TabUnstyled>TEAL</TabUnstyled>
                        <TabUnstyled>Base64</TabUnstyled>
                      </TabsListUnstyled>
                      <TabPanelUnstyled value={0}>
                        <Copyable copyableText={disassembledClearStateProgram}>
                          <pre className={`${styles["teal-box"]} hljs`}>
                            <code
                              className="language-lua"
                              dangerouslySetInnerHTML={{
                                __html: hljs.highlight(
                                  disassembledClearStateProgram,
                                  {
                                    language: "lua",
                                  }
                                ).value,
                              }}
                            />
                          </pre>
                        </Copyable>
                      </TabPanelUnstyled>
                      <TabPanelUnstyled value={1}>
                        <Copyable copyableText={appTx["clear-state-program"]} />
                      </TabPanelUnstyled>
                    </TabsUnstyled>
                  ) : (
                    <Copyable copyableText={appTx["clear-state-program"]} />
                  )}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
export default ApplicationTransactionInfo;
