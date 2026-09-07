#!/usr/bin/env python3
"""Generate comprehensive RH SPARS seed SQL from Final integrated tool structure."""

from pathlib import Path


def q(text, weight=1, order=1):
    return (text.replace("'", "''"), weight, order)


domains = []

# ========== MATERNITY ==========
mat_areas = []
mat_areas.append(
    (
        "Health Information System - Maternity Register",
        1,
        [
            q("Integrated Maternity Register is available", 1, 1),
            q("Ten random entries in the Integrated Maternity Register are completely filled", 1, 2),
            q("Monthly summary totals are present in the Integrated Maternity Register", 1, 3),
        ],
    )
)
mat_areas.append(
    (
        "Infrastructure - Maternity confidentiality and space",
        2,
        [
            q("Maternity: space for this service exists", 1, 1),
            q("Maternity: space/room ensures privacy for every client", 1, 2),
            q("Maternity: safe place to keep case notes so non-health professionals cannot read them", 1, 3),
            q("Facility has separate arrangements or waiting areas for young people", 1, 4),
            q("Clean toilet facilities clearly labelled male or female", 1, 5),
            q("Facility has a functional placenta pit", 1, 6),
        ],
    )
)
mat_areas.append(
    (
        "Job aids and protocols - Maternity unit",
        3,
        [
            q("Essential maternal and newborn care guidelines available at maternity workstation", 1, 1),
            q("PPH protocol/job aid available", 1, 2),
            q("Hypertension in pregnancy protocol/job aid available", 1, 3),
            q("PPFP compendium available", 1, 4),
            q("PAC FP compendium/protocols available", 1, 5),
            q("HIV EMTCT guidelines available", 1, 6),
            q("Referral forms available", 1, 7),
            q("Consent forms available", 1, 8),
            q("Partograph available", 1, 9),
        ],
    )
)
mat_areas.append(
    (
        "Infection prevention practices - Maternity unit",
        4,
        [
            q(
                "Colour-coded bins all present: non-infectious (black), infectious (yellow), "
                "highly infectious (red), pharmaceutical (brown), domestic (green)",
                1,
                1,
            ),
            q("Sharps boxes available", 1, 2),
            q("Handwashing facilities available", 1, 3),
            q("Sterilizer/autoclave available and functional for reusable equipment", 1, 4),
        ],
    )
)
equip = [
    "Neonatal ambu bag and mask",
    "Weighing scale for newborns",
    "Penguin sucker",
    "Neonatal resuscitation table",
    "Sterilization equipment",
    "Vacuum aspirator or D and C kit",
    "Suction apparatus",
    "Manual vacuum extractor",
    "Delivery bed",
    "Delivery bed for disabled",
    "Delivery kit",
    "Examination light",
    "Emergency transport",
    "PPH kit",
    "PET kit",
    "Fridge or cold box/vaccine carrier with cold ice packs for storage of oxytocin",
]
mat_areas.append(
    (
        "Essential equipment availability - Maternity",
        5,
        [q(f"{e} available and functional", 1, i + 1) for i, e in enumerate(equip)],
    )
)
meds = [
    "Magnesium sulphate inj. (ampoule)",
    "Calcium gluconate inj. (ampoule)",
    "Dextrose 50% (bottle)",
    "Hydrocortisone inj. (vial)",
    "Misoprostol (tablet)",
    "Vitamin K",
    "Oxytocin",
    "Tranexamic acid",
    "Heat-stable carbetocin",
]
mat_areas.append(
    (
        "Emergency medicines - Maternity",
        6,
        [q(f"{m} available", 1, i + 1) for i, m in enumerate(meds)],
    )
)
mat_areas.append(
    (
        "Human resources - Staff training (Maternity)",
        7,
        [
            q(
                "Any staff received training in maternity (on-job or workshop) within the past year "
                "(evidence: training register, CME book, and/or MoH iHRIS)",
                1,
                1,
            ),
        ],
    )
)
mat_areas.append(
    (
        "Functionality of committees - Department meeting (Maternity)",
        8,
        [
            q("Active department committee with regular meetings (last 2 meetings)", 1, 1),
            q("Data issues discussed in these meetings", 1, 2),
            q("RH supplies and commodities discussed in these meetings", 1, 3),
        ],
    )
)
mat_areas.append(
    (
        "Support supervision - Maternity",
        9,
        [
            q(
                "Maternity supervised by facility in-charge or senior management at least once "
                "during the previous quarter",
                1,
                1,
            ),
        ],
    )
)
domains.append(("maternity", "Maternity services", 1, mat_areas))

# ========== FP ==========
fp_areas = []
fp_areas.append(
    (
        "Health Information System - FP Register",
        1,
        [
            q("Integrated FP Register is available", 1, 1),
            q("Ten random entries in the Integrated FP Register are completely filled", 1, 2),
            q("Monthly summary totals are present in the Integrated FP Register", 1, 3),
        ],
    )
)
fp_areas.append(
    (
        "Infrastructure - Family planning confidentiality and space",
        2,
        [
            q("Family planning: space for this service exists", 1, 1),
            q("Family planning: consultation room ensures privacy for every client", 1, 2),
            q(
                "Family planning: safe place to keep case notes so non-health professionals cannot read them",
                1,
                3,
            ),
            q("Facility has separate arrangements or waiting areas for young people", 1, 4),
            q("Clean toilet facilities clearly labelled male or female", 1, 5),
            q(
                "Facilities for waste disposal (dump pit, incinerator) or disposal plan available",
                1,
                6,
            ),
        ],
    )
)
fp_aids = [
    "Method mix chart",
    "Myths and misconceptions pocket book",
    "FP counselling flip chart",
    "PPFP compendium",
    "PAC FP compendium/protocols",
    "MEC wheel",
    "Demonstration models",
    "Referral forms",
    "Consent forms",
]
fp_areas.append(
    (
        "Job aids and protocols - Family planning space(s)",
        3,
        [q(f"{a} available at FP workstation", 1, i + 1) for i, a in enumerate(fp_aids)],
    )
)
fp_areas.append(
    (
        "Infection prevention practices - Family planning space(s)",
        4,
        [
            q(
                "Colour-coded bins all present: non-infectious (black), infectious (yellow), "
                "highly infectious (red), pharmaceutical (brown), domestic (green)",
                1,
                1,
            ),
            q("Sharps boxes available", 1, 2),
            q("Handwashing facilities available", 1, 3),
            q("Sterilizer/autoclave available and functional for reusable equipment", 1, 4),
        ],
    )
)
fp_equip_main = [
    "Weighing scale available and functional",
    "BP machine available and functional",
    "Implant insertion/removal kit complete (mosquito artery forceps curved and straight; "
    "small dissecting forceps with teeth; scissors small sharp; scalpel handle with blade No. 11 or 15)",
    "IUD insertion/removal kit complete (Cusco speculum; tenaculum; uterine sound; "
    "sponge-holding forceps; long curved or ring forceps; scissors)",
    "BTL kit available and functional",
    "Vasectomy kit available and functional",
    "Angle light / lamp available and functional",
    "Cervical cancer screening and preventive treatment available",
    "Breast examination service available",
]
fp_areas.append(
    (
        "Essential equipment availability - Family planning",
        5,
        [q(e, 1, i + 1) for i, e in enumerate(fp_equip_main)],
    )
)
fp_areas.append(
    (
        "Human resources - Staff training (Family planning)",
        6,
        [
            q(
                "Any staff received training in family planning within the past year "
                "(evidence: training register, CME book, and/or MoH iHRIS)",
                1,
                1,
            ),
            q("Staff undergone proficiency testing for family planning", 1, 2),
        ],
    )
)
fp_areas.append(
    (
        "Functionality of committees - Department meeting (FP)",
        7,
        [
            q("Active department committee with regular meetings (last 2 meetings)", 1, 1),
            q("Data issues discussed in these meetings", 1, 2),
            q("RH supplies and commodities discussed in these meetings", 1, 3),
        ],
    )
)
fp_areas.append(
    (
        "Support supervision - Family planning",
        8,
        [
            q(
                "Family planning supervised by facility in-charge or senior management at least once "
                "during the previous quarter",
                1,
                1,
            ),
        ],
    )
)
domains.append(("fp", "Family planning services", 2, fp_areas))

# ========== ANC ==========
anc_areas = []
anc_areas.append(
    (
        "Health Information System - ANC Register",
        1,
        [
            q("Integrated ANC Register is available", 1, 1),
            q("Ten random entries in the Integrated ANC Register are completely filled", 1, 2),
            q("Monthly summary totals are present in the Integrated ANC Register", 1, 3),
        ],
    )
)
anc_areas.append(
    (
        "Infrastructure - Antenatal clinic confidentiality and space",
        2,
        [
            q("Antenatal clinic: space for this service exists", 1, 1),
            q("Antenatal clinic: consultation room ensures privacy for every client", 1, 2),
            q(
                "Antenatal clinic: safe place to keep case notes so non-health professionals cannot read them",
                1,
                3,
            ),
            q("Facility has separate arrangements or waiting areas for young people", 1, 4),
            q("Clean toilet facilities clearly labelled male or female", 1, 5),
            q(
                "Facilities for waste disposal (dump pit, incinerator) or disposal plan available",
                1,
                6,
            ),
        ],
    )
)
anc_aids = [
    "Essential Maternal and Newborn Care guidelines",
    "STI screening protocols",
    "HIV EMTCT guidelines",
    "PPH client-facing posters/messages (high-risk conditions)",
    "PPFP compendium",
    "Referral forms",
    "Consent forms",
]
anc_areas.append(
    (
        "Job aids and protocols - Antenatal unit",
        3,
        [q(f"{a} available at ANC workstation", 1, i + 1) for i, a in enumerate(anc_aids)],
    )
)
anc_areas.append(
    (
        "Infection prevention practices - Antenatal clinic",
        4,
        [
            q(
                "Colour-coded bins all present: non-infectious (black), infectious (yellow), "
                "highly infectious (red), pharmaceutical (brown), domestic (green)",
                1,
                1,
            ),
            q("Sharps boxes available", 1, 2),
            q("Handwashing facilities available", 1, 3),
            q("Sterilizer/autoclave available and functional for reusable equipment", 1, 4),
        ],
    )
)
anc_equip = ["BP machine", "Fetoscope", "Ultrasound", "MUAC tape", "Weighing scale", "Height meter"]
anc_areas.append(
    (
        "Essential equipment availability - ANC",
        5,
        [q(f"{e} available and in use", 1, i + 1) for i, e in enumerate(anc_equip)],
    )
)
anc_areas.append(
    (
        "Human resources - Staff training (ANC)",
        6,
        [
            q(
                "Any staff received training in ANC within the past year "
                "(evidence: training register, CME book, and/or MoH iHRIS)",
                1,
                1,
            ),
        ],
    )
)
anc_areas.append(
    (
        "Functionality of committees - Department meeting (ANC)",
        7,
        [
            q("Active department committee with regular meetings (last 2 meetings)", 1, 1),
            q("Data issues discussed in these meetings", 1, 2),
            q("RH supplies and commodities discussed in these meetings", 1, 3),
        ],
    )
)
anc_areas.append(
    (
        "Support supervision - ANC services",
        8,
        [
            q(
                "ANC services supervised by facility in-charge or senior management at least once "
                "during the previous quarter",
                1,
                1,
            ),
        ],
    )
)
domains.append(("anc", "Antenatal care services", 3, anc_areas))

# ========== COMMODITIES ==========
commodities = [
    "Medroxyprogesterone Acetate (Depo-Provera) 150mg/ml for I/M injection",
    "Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection",
    "Etonogestrel 68mg Implant (Implanon)",
    "Levonorgestrel 2x75mg Implant (Levoplant)",
    "Levonorgestrel 2x75mg Implant (Jadelle)",
    "Levonorgestrel 0.15mg + Ethinyl Estradiol 0.03mg Tablets (28 Tablets)",
    "Levonorgestrel 0.75mg, 2 Tablets (ECP)",
    "Hormonal IUD",
    "Cu-T Intrauterine Device",
    "Male Condom",
    "Sulfadoxine 500 + Pyrimethamine 25mg Tablets",
    "Ferrous Sulphate/Fumarate + Folic Acid Tablets",
    "Mebendazole 500mg Tablets",
    "Oxytocin Injection 10 i.u./mL",
    "Magnesium Sulphate Injection 50%",
    "Female condoms",
    "Moon beads",
]
com_areas = []
com_areas.append(
    (
        "Stock management - Availability of items (B1)",
        1,
        [q(f"{c}: item available on day of survey", 1, i + 1) for i, c in enumerate(commodities)],
    )
)
com_areas.append(
    (
        "Stock management - No expired stock (B)",
        2,
        [
            q(f"{c}: no expired quantity in stock (score 1 if NO expired stock present)", 1, i + 1)
            for i, c in enumerate(commodities)
        ],
    )
)
com_areas.append(
    (
        "Stock management - Availability of stock cards (B2)",
        3,
        [q(f"{c}: stock card available", 1, i + 1) for i, c in enumerate(commodities)],
    )
)
sample_stock = commodities[:5]
com_areas.append(
    (
        "Stock management - Physical count last 3 months (B3)",
        4,
        [
            q(
                f"{c}: physical count done for last three months and PC marked in stock card",
                1,
                i + 1,
            )
            for i, c in enumerate(sample_stock)
        ],
    )
)
com_areas.append(
    (
        "Stock management - Correct filling of stock cards (B4)",
        5,
        [
            q(f"{c}: stock card filled correctly with name, strength, dosage form, AMC", 1, i + 1)
            for i, c in enumerate(sample_stock)
        ],
    )
)
com_areas.append(
    (
        "Stock management - Updating/accuracy of stock card balance (B5)",
        6,
        [
            q(f"{c}: stock card balance and physical count agree 100%", 1, i + 1)
            for i, c in enumerate(sample_stock)
        ],
    )
)
com_areas.append(
    (
        "Stock management - Accuracy of AMC (B6)",
        7,
        [
            q(f"{c}: calculated AMC same as recorded AMC within +/- 10%", 1, i + 1)
            for i, c in enumerate(sample_stock)
        ],
    )
)
# Column P: Stock Status — Optimum=1, Overstock/Understock=0
com_areas.append(
    (
        "Stock management - Stock status optimum (P)",
        8,
        [
            q(
                f"{c}: stock status is Optimum for the next 3 months "
                "(Yes=Optimum; No=Overstock or Understock)",
                1,
                i + 1,
            )
            for i, c in enumerate(commodities)
        ],
    )
)
com_areas.append(
    (
        "Procurement planning and ordering - Awareness",
        9,
        [
            q(
                "Facility undertook joint planning (e.g., with IP) and staff are aware of planned "
                "quarterly/monthly Family Planning Outreaches",
                1,
                1,
            ),
            q(
                "Facility staff aware of emergency ordering or redistribution procedures for RH "
                "commodities (notify DHT of shortages or excesses)",
                1,
                2,
            ),
            q("Facility is using the online stock status facility system", 1, 3),
        ],
    )
)
proc_qs = [
    q(
        "EMHS Procurement Plan for the current financial year available at the health facility",
        1,
        1,
    )
]
for i, c in enumerate(commodities):
    proc_qs.append(q(f"Procurement plan includes: {c}", 1, i + 2))
com_areas.append(("Procurement plan contents", 10, proc_qs))
com_areas.append(
    (
        "Order integration",
        11,
        [
            q(
                "Facility submitted RH commodity orders alongside EMHS, Lab, ART, and TB orders",
                1,
                1,
            ),
        ],
    )
)
com_areas.append(
    (
        "Order timeliness",
        12,
        [
            q(
                "Most recent RH commodity order was submitted timely (before NMS/JMS order deadline)",
                1,
                1,
            ),
        ],
    )
)
sample5 = [
    "Medroxyprogesterone Acetate (Sayana Press) 104mg/0.65ml for S/C injection",
    "Etonogestrel 68mg Implant (Implanon)",
    "Cu-T Intrauterine Device",
    "Ferrous Sulphate/Fumarate + Folic Acid Tablets",
    "Magnesium Sulphate Injection 50%",
]
oq = []
order = 1
for c in sample5:
    oq.append(q(f"Order quality - {c}: opening balances on both order forms agree", 1, order))
    order += 1
    oq.append(q(f"Order quality - {c}: number of clients recorded on last RH order form", 1, order))
    order += 1
    oq.append(
        q(
            f"Order quality - {c}: client counts on RH order form vs FP register agree within +/- 10%",
            1,
            order,
        )
    )
    order += 1
com_areas.append(("Order quality - last supplied order cycle", 13, oq))
domains.append(("commodities", "Drugs and supplies", 4, com_areas))

# ========== HEALTH INFORMATION ==========
hi_areas = []
tools = [
    "Family Planning Register",
    "Integrated ANC Register",
    "Integrated Maternity Register",
    "Postnatal Register",
    "Immunisation child register",
    "Monthly summaries report forms 105",
    "Community service registers and Summary Report 097b or eCHIS",
    "Dispensing log",
    "Electronic Medical Records EMR (e.g., eAFYA, Clinic Master, Ug-EMR)",
]
hi_qs = []
order = 1
for t in tools:
    hi_qs.append(q(f"{t}: available", 1, order))
    order += 1
    hi_qs.append(q(f"{t}: ten random entries completely filled (where applicable)", 1, order))
    order += 1
    hi_qs.append(q(f"{t}: monthly summary totals present (where applicable)", 1, order))
    order += 1
hi_areas.append(("Availability and utilisation of HMIS tools", 1, hi_qs))
hi_areas.append(
    (
        "Reporting - Completeness of HMIS 105",
        2,
        [
            q("Facility submitted reports for the two months prior to the supervision", 1, 1),
            q("Facility submitted information for ANC services in HMIS 105 Section 2.1", 1, 2),
            q("Facility submitted information for Maternity services in HMIS 105 Section 2.2", 1, 3),
            q(
                "Facility submitted information for Family Planning services in HMIS 105 Section 2.4",
                1,
                4,
            ),
            q(
                "Facility submitted information about dispensing of Family Planning methods in "
                "HMIS 105 Section 2.4.2",
                1,
                5,
            ),
            q(
                "Facility submitted all data for all these commodities in HMIS 105 Section 6 "
                "(qty consumed, days out of stock, stock on hand, qty expired for: DMPA, SP tablets, "
                "Misoprostol 200mcg, Implanon, Oxytocin, Chlorhexidine Gel, Mama Kits)",
                1,
                6,
            ),
        ],
    )
)
accuracy = [
    "Number of Pregnant Women receiving at least 30 tablets of Folic Acid and Iron Sulphate at ANC 1st contact/visit (AN35)",
    "Number of pregnant women dewormed or receiving Mebendazole (AN34)",
    "Total number of deliveries (MA04)",
    "Total number (new + revisits) of Injectable DMPA (IM e.g., Depo) users (FP06)",
    "Total number (new + revisits) of users of 3 year implant e.g., Implanon NXT, Levoplant (FP10)",
    "Total number (new + revisits) of users of IUD-Copper-T (FP13)",
    "Quantity of DMPA issued in the store (SS02)",
    "Quantity of Sulfadoxine/Pyrimethamine tablets issued in the store (SS04)",
    "Quantity of Oxytocin issued in the store (SS31)",
]
acc_qs = []
order = 1
for a in accuracy:
    acc_qs.append(q(f"HMIS accuracy - {a}: information available from the last report", 1, order))
    order += 1
    acc_qs.append(q(f"HMIS accuracy - {a}: data agree or differ by no more than +/- 10%", 1, order))
    order += 1
hi_areas.append(("Accuracy of HMIS 105 Report", 3, acc_qs))
domains.append(("health_information", "Health Information", 5, hi_areas))

# ========== GENERAL FACILITY ==========
gf_areas = []
services = ["Family planning services", "Antenatal services", "Maternity unit"]
sign_qs = []
order = 1
for s in services:
    sign_qs.append(
        q(
            f"Signage - {s}: poster clearly shows availability of these services (English and local language)",
            1,
            order,
        )
    )
    order += 1
    sign_qs.append(q(f"Signage - {s}: clearly shows the day and time these services are available", 1, order))
    order += 1
    sign_qs.append(
        q(
            f"Signage - {s}: services provided at least 5 days a week for at least 6 hours "
            "(maternity should be 24 hours)",
            1,
            order,
        )
    )
    order += 1
    sign_qs.append(q(f"Signage - {s}: clear within-facility directional signage to the clinic", 1, order))
    order += 1
gf_areas.append(("Signage / poster", 1, sign_qs))
committees = [
    "Health unit management committee / Hospital Management board",
    "Senior management meeting",
    "QI meeting",
]
com_qs = []
order = 1
for c in committees:
    com_qs.append(q(f"{c}: active committee with regular meetings (last 2 meetings)", 1, order))
    order += 1
    com_qs.append(q(f"{c}: MCH issues discussed in these meetings", 1, order))
    order += 1
    com_qs.append(q(f"{c}: RH supplies and commodities discussed", 1, order))
    order += 1
gf_areas.append(("Functionality of facility committees", 2, com_qs))
cadres = [
    "Obstetrician and gynaecologists",
    "Medical officers",
    "Midwives",
    "Nurses",
    "Anaesthetic officers",
    "Assistant inventory management officer",
    "Health information staff",
    "Laboratory staff",
    "Radiography / radiology staff",
]
gf_areas.append(
    (
        "Human resources - Staffing levels",
        3,
        [
            q(
                f"{c}: staffing level at or above 80% of recommended norm (per duty roster)",
                1,
                i + 1,
            )
            for i, c in enumerate(cadres)
        ],
    )
)
bsc = cadres[:6]
gf_areas.append(
    (
        "Staff performance management - Balanced score cards",
        4,
        [
            q(
                f"{c}: sampled balanced score cards include clear MCH indicators (FP, ANC, Maternity)",
                1,
                i + 1,
            )
            for i, c in enumerate(bsc)
        ],
    )
)
domains.append(("general_facility", "General facility", 6, gf_areas))


def main():
    lines = [
        "-- RH SPARS (Integrated Reproductive Health Support Supervision Tool) seed",
        "-- Generated from Final integrated tool.xlsx structure",
        "-- Scoring: Yes=1, No=0, NA excluded from numerator and denominator",
        "-- section% = sum/(n-NA)*100; domain% = avg(section%); grand% = avg(domain%)",
        "",
        "BEGIN;",
        "",
    ]

    for code, name, dorder, areas in domains:
        name_esc = name.replace("'", "''")
        lines.append(
            "INSERT INTO rhspars_domains (code, name, display_order) "
            f"SELECT '{code}', '{name_esc}', {dorder} "
            f"WHERE NOT EXISTS (SELECT 1 FROM rhspars_domains WHERE code='{code}');"
        )
        lines.append("")
        for aname, aorder, questions in areas:
            aname_esc = aname.replace("'", "''")
            lines.append(
                "INSERT INTO rhspars_thematic_areas (domain_id, name, display_order) "
                f"SELECT d.id, '{aname_esc}', {aorder} FROM rhspars_domains d "
                f"WHERE d.code='{code}' AND NOT EXISTS ("
                "SELECT 1 FROM rhspars_thematic_areas ta "
                f"WHERE ta.domain_id=d.id AND ta.name='{aname_esc}');"
            )
            lines.append("")
            lines.append(
                "INSERT INTO rhspars_questions (thematic_area_id, question_text, score_weight, display_order)"
            )
            lines.append("SELECT ta.id, q.question_text, q.score_weight, q.display_order")
            lines.append("FROM rhspars_thematic_areas ta")
            lines.append("JOIN rhspars_domains d ON ta.domain_id = d.id,")
            lines.append("(VALUES")
            vals = [f"    ('{text}', {w}, {o})" for text, w, o in questions]
            lines.append(",\n".join(vals))
            lines.append(") AS q(question_text, score_weight, display_order)")
            lines.append(f"WHERE d.code = '{code}' AND ta.name = '{aname_esc}'")
            lines.append(
                "AND NOT EXISTS (SELECT 1 FROM rhspars_questions x "
                "WHERE x.thematic_area_id = ta.id AND x.question_text = q.question_text);"
            )
            lines.append("")

    lines.append("COMMIT;")
    lines.append("")
    total_q = sum(len(qs) for _, _, _, areas in domains for _, _, qs in areas)
    total_a = sum(len(areas) for _, _, _, areas in domains)
    lines.append(f"-- Seeded {len(domains)} domains, {total_a} thematic areas, {total_q} questions")

    out = Path(__file__).resolve().parents[1] / "seed-rhspars-questions.sql"
    out.write_text("\n".join(lines), encoding="utf-8")
    print(f"Wrote {out}")
    print(f"Domains={len(domains)} areas={total_a} questions={total_q}")


if __name__ == "__main__":
    main()
