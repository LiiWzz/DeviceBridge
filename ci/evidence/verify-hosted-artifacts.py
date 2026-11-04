# Executed evidence verifier for PR #18, run 41911387490.
# argv[1] is connector-provided artifact metadata plus short-lived download URLs.
# Download URLs are never included in the saved verification records.
import sys,json,urllib.request,hashlib,zipfile,io,pathlib
items=json.loads(sys.argv[1])
out=pathlib.Path("dist/ci/downloads");out.mkdir(parents=True,exist_ok=True)
results=[]
for p in items:
    data=urllib.request.urlopen(p["url"],timeout=30).read()
    assert "sha256:"+hashlib.sha256(data).hexdigest()==p["digest"],"ZIP digest mismatch"
    (out/(str(p["id"])+".zip")).write_bytes(data)
    with zipfile.ZipFile(io.BytesIO(data)) as z:
        source_paths=[name for name in z.namelist() if name=="SOURCE.txt" or name.endswith("/SOURCE.txt")]
        assert len(source_paths)==1,source_paths
        prefix=source_paths[0][:-len("SOURCE.txt")]
        source=dict(line.split("=",1) for line in z.read(prefix+"SOURCE.txt").decode().splitlines())
        assert source["workspace_dirty"]=="false",source
        assert source["checked_out_commit"]=="f91559e3d0290b12a8cdb8ebf1ce53f6d47156e6",source
        assert source["run_id"]=="41911387490",source
        checked=[]
        if source["job"]!="workflow":
            for line in z.read(prefix+"SHA256SUMS").decode().splitlines():
                digest,name=line.split("  ",1)
                assert hashlib.sha256(z.read(prefix+name)).hexdigest()==digest,name
                checked.append(name)
        else:
            assert set(z.namelist())=={"SOURCE.txt","commands.txt"}
        result={"artifact_id":p["id"],"name":p["name"],"zip_sha256":p["digest"],"zip_digest_valid":True,"manifest_files_verified":checked,"source":source}
        results.append(result)
        (out/(str(p["id"])+"-verification.json")).write_text(json.dumps(result,indent=2)+"\n")
        print(source["job"]+": ZIP digest, source provenance and "+str(len(checked))+" manifest entries verified",flush=True)
