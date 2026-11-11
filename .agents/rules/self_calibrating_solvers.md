# Sovereign Rule: Self-Calibrating Solvers & Zero Hardcoded Empirical Constants

> **Scope:** Applicable to all physical, statistical, numerical, forecasting, and machine-learning solvers across the sovereign fleet.  
> **Governed Workspaces:** All `00flow` and `00xper` computational, scientific, and analytics workspaces.

---

<hard_invariants>
## 1. Strict Epistemic Separation of Constants vs. Calibrated Parameters

1. **PERMITTED UNIVERSAL CONSTANTS:** Only universal physical, astronomical, or exact mathematical constants may be written as literal constants in solver code:
   - Standard physical constants (e.g., standard gravity $g = 9.80665\text{ m/s}^2$, standard seawater density $\rho_0 = 1025.0\text{ kg/m}^3$, standard atmospheric pressure $P_0 = 1013.25\text{ hPa}$).
   - Exact mathematical, spectral, or astronomical identities (e.g., Doodson tidal speeds, Rayleigh spectral identity $H_{m0} = 4.004\sqrt{m_0}$, RK4 Butcher tableau weights).
   - Hardware/sensor specification constants defined by physical instrumentation (e.g., NOAA NGWLMS $181\text{-second}$ 1-Hz averaging window).
2. **PROHIBITED EMPIRICAL MAGIC NUMBERS:** Agents are strictly forbidden from hardcoding empirical, statistical, or tuning constants in solvers, including:
   - Static reference standard deviations or variance coupling ratios in heteroscedastic GLS/Cholesky solvers.
   - Arbitrary Tikhonov/ridge regularization multipliers ($\lambda$) or step-function shrinkage thresholds.
   - Fixed Kalman, Holt-Winters, or kinematic filter gains ($\alpha, \beta, K$).
   - Fixed ensemble perturbation percentages or static climatological reference baselines (e.g., hardcoded reference temperatures, pressures, or irradiances).
3. **MANDATORY CAUSAL SELF-CALIBRATION ($t \le t_0$):** Every empirical or regime-dependent parameter must be:
   - Estimated dynamically from causal spin-up observations ($t \le t_0$) via Empirical Bayes, Restricted Maximum Likelihood (REML), Gauss-Markov variance decomposition, Hoerl-Kennard-Baldwin noise-to-signal ratios, Yule-Walker/Riccati state-space equations, or sample moments; OR
   - Derived in closed form from physical conservation laws and station/domain geometry.
4. **MANDATORY NON-DEGENERACY VERIFICATION:** Any solver introducing or modifying calibrated parameters must include a unit test in the workspace test suite (`realize.exe`) asserting that the learned parameters are positive/well-conditioned and take distinct values across distinct stations, regimes, or datasets.
</hard_invariants>
